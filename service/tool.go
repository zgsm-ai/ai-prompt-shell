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
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type CallStats struct {
	Count    int64
	Duration time.Duration
	LastErr  error
}

var registry = dao.NewToolCache()
var callStats = sync.Map{}

func GetTool(toolId string) (dao.Tool, error) {
	t, exists := registry.Get(toolId)
	if !exists {
		return dao.Tool{}, os.ErrNotExist
	}
	return t, nil
}

func ListTool() map[string]dao.Tool {
	return registry.All()
}

func ToolIDs() ([]string, error) {
	var results []string
	for k, _ := range registry.All() {
		results = append(results, k)
	}
	return results, nil
}

func Call(ctx context.Context, toolId string, args map[string]interface{}) (interface{}, error) {
	tool, ok := registry.Get(toolId)
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", toolId)
	}
	if err := utils.ValidateVariables(args, tool.Parameters); err != nil {
		return nil, err
	}
	return callTool(ctx, tool, args)
}

func callTool(ctx context.Context, tool dao.Tool, args map[string]interface{}) (interface{}, error) {
	// 根据工具类型调用不同的执行逻辑
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

func callRestfulTool(ctx context.Context, tool dao.Tool, params map[string]interface{}) (interface{}, error) {
	// 2. URL验证
	if tool.URL == "" {
		return nil, fmt.Errorf("missing URL for tool %s", tool.Name)
	}
	if _, err := http.NewRequest("GET", tool.URL, nil); err != nil {
		return nil, fmt.Errorf("invalid URL %q for tool %s: %v", tool.URL, tool.Name, err)
	}

	// 3. 构建请求和重试逻辑
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var lastErr error
	var result interface{}

	// 最多重试3次
	for i := 0; i < 3; i++ {
		reqBody, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %v", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", tool.URL, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "shenma-tool-caller/1.0")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d): %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}

		defer resp.Body.Close()

		// 4. 处理响应状态码
		if resp.StatusCode >= 500 {
			// 服务器错误，尝试重试
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			lastErr = fmt.Errorf("server error (status %d): %s", resp.StatusCode, string(body))
			continue
		}
		if resp.StatusCode >= 400 {
			// 客户端错误，立即返回
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			return nil, fmt.Errorf("tool call failed (status %d): %s", resp.StatusCode, string(body))
		}

		// 5. 解析成功响应
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %v", err)
		}

		return result, nil
	}

	// 6. 超过最大重试次数
	return nil, fmt.Errorf("max retries reached for tool %s: %v", tool.Name, lastErr)
}

func callGRPCTool(ctx context.Context, tool dao.Tool, params map[string]interface{}) (interface{}, error) {
	// 2. URL验证
	if tool.URL == "" {
		return nil, fmt.Errorf("missing gRPC endpoint URL for tool %s", tool.Name)
	}

	var lastErr error
	var result interface{}
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 3. 建立连接 (带超时控制)
		connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(connCtx, tool.URL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())
		if err != nil {
			lastErr = fmt.Errorf("failed to connect to gRPC server (attempt %d): %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
			continue
		}
		defer conn.Close()

		// 4. 准备请求数据
		reqBody, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %v", err)
		}

		resp := new(bytes.Buffer)
		method := fmt.Sprintf("%s/%s", tool.Module, tool.Name)

		// 5. 执行gRPC调用 (带超时控制)
		callCtx, callCancel := context.WithTimeout(ctx, 10*time.Second)
		defer callCancel()

		err = conn.Invoke(callCtx, method, reqBody, resp)
		if err != nil {
			lastErr = fmt.Errorf("gRPC call failed (attempt %d): %v", i+1, err)
			if isRetriableError(err) {
				time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
				continue
			}
			return nil, lastErr
		}

		// 6. 解析响应
		if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
			lastErr = fmt.Errorf("failed to parse gRPC response: %v", err)
			break // 解析错误不重试
		}

		// 7. 更新统计信息
		stats, _ := callStats.LoadOrStore(tool.Name, &CallStats{})
		stats.(*CallStats).Count++
		return result, nil
	}

	// 8. 处理重试耗尽
	stats, _ := callStats.LoadOrStore(tool.Name, &CallStats{})
	stats.(*CallStats).LastErr = lastErr
	return nil, fmt.Errorf("max retries (%d) reached for gRPC tool %s: %v", maxRetries, tool.Name, lastErr)
}

func isRetriableError(err error) bool {
	// 可重试的错误包括: 网络超时、连接拒绝、服务不可用等
	return isDeadlineExceeded(err) || isUnavailable(err)
}

func isDeadlineExceeded(err error) bool {
	return err == context.DeadlineExceeded || status.Code(err) == codes.DeadlineExceeded
}

func isUnavailable(err error) bool {
	return status.Code(err) == codes.Unavailable
}

func callMCPTool(ctx context.Context, tool dao.Tool, params map[string]interface{}) (interface{}, error) {
	// 简单实现，实际项目中应该通过MCP客户端调用
	return nil, fmt.Errorf("MCP tool not implemented yet")
}
