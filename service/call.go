package service

import (
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CallStats struct {
	Count    int64
	Duration time.Duration
	LastErr  error
}

var callStats = sync.Map{}

/**
 * Execute tool call with validation
 * @param ctx Context for the call
 * @param t Tool definition
 * @param args Arguments for the tool
 * @return Execution result or error
 */
func Call(ctx context.Context, t *dao.Tool, args []interface{}) (interface{}, error) {
	if err := utils.ValidateArgs(args, t.Parameters); err != nil {
		return nil, err
	}
	return callTool(ctx, t, args)
}

/**
 * Route to appropriate tool executor based on tool type
 * @param ctx Context for the call
 * @param tool Tool definition
 * @param args Arguments for the tool
 * @return Execution result or error
 */
func callTool(ctx context.Context, tool *dao.Tool, args []interface{}) (interface{}, error) {
	switch tool.Type {
	case "restful":
		return callRestfulTool(ctx, tool, args)
	case "grpc":
		return callGRPCTool(ctx, tool, args)
	case "mcp":
		return callMCPTool(ctx, tool, args)
	default:
		return nil, fmt.Errorf("unsupported tool type: %s", tool.Type)
	}
}

/**
 * Call RESTful API tool
 * @param ctx Context for the call
 * @param tool RESTful tool definition
 * @param args Arguments for API call
 * @return API response or error
 */
func callRestfulTool(ctx context.Context, tool *dao.Tool, args []interface{}) (interface{}, error) {
	if tool.Restful == nil {
		return nil, fmt.Errorf("missing RESTful definition for tool %s", tool.Name)
	}
	if tool.Restful.Url == "" {
		return nil, fmt.Errorf("missing URL for tool %s", tool.Name)
	}
	if tool.Restful.Method == "" {
		return nil, fmt.Errorf("missing method for tool %s", tool.Name)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	reqBody, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, tool.Restful.Method, tool.Restful.Url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ai-prompt-shell/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tool %s request failed: %v", tool.Name, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("server error for tool %s (status %d): %s", tool.Name, resp.StatusCode, string(body))
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("tool %s call failed (status %d): %s", tool.Name, resp.StatusCode, string(body))
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %v", tool.Name, err)
	}

	return result, nil
}

/**
 * Call gRPC service tool without retry
 * @param ctx Context for the call
 * @param tool gRPC tool definition with endpoint and method
 * @param args Arguments for gRPC call
 * @return gRPC response or error
 */
func callGRPCTool(ctx context.Context, tool *dao.Tool, args []interface{}) (interface{}, error) {
	if tool.Grpc.Url == "" {
		return nil, fmt.Errorf("missing gRPC endpoint URL for tool %s", tool.Name)
	}

	// 1. Establish connection (with timeout)
	connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(connCtx, tool.Grpc.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// 2. Prepare request data
	reqBody, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp := new(bytes.Buffer)
	method := fmt.Sprintf("%s/%s", tool.Module, tool.Name)

	// 3. Execute gRPC call (with timeout)
	callCtx, callCancel := context.WithTimeout(ctx, 10*time.Second)
	defer callCancel()

	err = conn.Invoke(callCtx, method, reqBody, resp)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %v", err)
	}

	// 4. Parse response
	var result interface{}
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse gRPC response: %v", err)
	}

	// 5. Update statistics
	stats, _ := callStats.LoadOrStore(tool.Name, &CallStats{})
	stats.(*CallStats).Count++
	return result, nil
}

func callMCPTool(ctx context.Context, tool *dao.Tool, args []interface{}) (interface{}, error) {
	// Simple implementation, should use MCP client in production
	return nil, fmt.Errorf("MCP tool not implemented yet")
}
