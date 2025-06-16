package template

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func keyToJsonPath(key string) string {
	return strings.ReplaceAll(strings.TrimPrefix(key, "shenma:environs:"), ":", ".")
}

func keyToFunc(key string) string {
	return strings.ReplaceAll(strings.TrimPrefix(key, "shenma:tools:"), ":", "_")
}

type RenderStats struct {
	Total         int64
	Success       int64
	Failures      int64
	TotalDuration time.Duration
	AvgDuration   time.Duration
	CacheHits     int64
	CacheMisses   int64
	RecentErrors  []error
	mu            sync.RWMutex
}

type Manager struct {
	refreshInterval time.Duration
	templates       map[string]*template.Template
	environs        map[string]interface{}
	tools           map[string]interface{}
	funcMap         template.FuncMap
	stats           *RenderStats
}

func NewManager(refreshInterval time.Duration) *Manager {
	m := &Manager{
		refreshInterval: refreshInterval,
		templates:       make(map[string]*template.Template),
		environs:        make(map[string]interface{}),
		tools:           make(map[string]interface{}),
		funcMap:         make(template.FuncMap),
	}

	// 初始加载
	m.refreshAllData()
	go m.startRefreshTimer()
	return m
}

func (m *Manager) startRefreshTimer() {
	ticker := time.NewTicker(m.refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.refreshAllData()
	}
}

func (m *Manager) refreshAllData() {
	m.refreshTemplates()
	m.refreshEnvirons()
	m.refreshTools()
}

// callMCP 调用MCP服务的简化实现
func callMCP(req map[string]interface{}) (interface{}, error) {
	// 这个简化实现直接返回空值
	// 实际项目中应该通过MCP客户端进行调用
	return nil, nil
}

func (m *Manager) LoadTemplates() (map[string]string, error) {
	vals, err := dao.LoadJsons("shenma:templates:")
	if err != nil {
		return nil, err
	}
	results := make(map[string]string)
	for key, val := range vals {
		results[key] = val.(string)
	}
	return results, nil
}

func (m *Manager) LoadEnvirons() (map[string]interface{}, error) {
	vals, err := dao.LoadJsons("shenma:environs:")
	if err != nil {
		return nil, err
	}
	results := make(map[string]interface{})
	for key, val := range vals {
		results[key] = val
	}
	return results, nil
}

func (m *Manager) LoadTools() (map[string]interface{}, error) {
	vals, err := dao.LoadJsons("shenma:tools:")
	if err != nil {
		return nil, err
	}
	results := make(map[string]interface{})
	for key, val := range vals {
		results[key] = val
	}
	return results, nil
}

func (m *Manager) refreshTemplates() {
	templates, err := m.LoadTemplates()
	if err != nil {
		return
	}

	for key, content := range templates {
		t, err := template.New(key).Funcs(m.funcMap).Parse(content)
		if err != nil {
			continue
		}
		m.templates[key] = t
	}
}

func (m *Manager) refreshEnvirons() {
	environs, err := m.LoadEnvirons()
	if err != nil {
		return
	}

	m.environs = environs
}

func (m *Manager) refreshTools() {
	tools, err := m.LoadTools()
	if err != nil {
		return
	}

	m.tools = tools
	// 更新模板函数表
	m.updateFuncMap()
}

func (m *Manager) updateFuncMap() {
	for name, tool := range m.tools {
		// 解析工具定义
		toolMap, ok := tool.(map[string]interface{})
		if !ok {
			continue
		}

		m.funcMap[name] = func(args ...interface{}) (interface{}, error) {
			// 1. 参数验证
			if len(args) < 1 {
				return nil, utils.ErrInvalidVariable
			}

			// 2. 根据工具类型调用不同的执行逻辑
			switch toolMap["type"] {
			case "restful":
				// TODO: 实现RESTful工具调用
				return callRestfulTool(toolMap, args)
			case "grpc":
				// TODO: 实现GRPC工具调用
				return callGRPCTool(toolMap, args)
			case "mcp":
				// TODO: 实现MCP工具调用
				return callMCPTool(toolMap, args)
			default:
				return nil, utils.ErrToolCallFailed
			}
		}
	}
}

// callRestfulTool 调用RESTful工具
func callRestfulTool(tool map[string]interface{}, args []interface{}) (interface{}, error) {
	// 1. 解析工具定义
	url, ok := tool["url"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}
	method, _ := tool["method"].(string)
	if method == "" {
		method = "POST"
	}

	// 2. 构建请求体
	var body interface{}
	if len(args) > 0 {
		body = args[0]
	}

	// 3. 发送HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. 处理响应
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("tool call failed with status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// callGRPCTool 调用GRPC工具
func callGRPCTool(tool map[string]interface{}, args []interface{}) (interface{}, error) {
	// 1. 解析工具定义
	service, ok := tool["service"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}
	method, ok := tool["method"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}
	address, ok := tool["address"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}

	// 2. 建立GRPC连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	// 3. 准备请求参数
	reqBytes, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	// 4. 反射调用GRPC方法(简化实现)
	ctx := context.Background()
	res, err := invokeGRPCMethod(ctx, conn, service, method, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("grpc call %s/%s failed: %v", service, method, err)
	}

	return res, nil
}

// callMCPTool 调用MCP工具
func callMCPTool(tool map[string]interface{}, args []interface{}) (interface{}, error) {
	// 1. 解析工具定义
	server, ok := tool["server"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}
	toolName, ok := tool["tool"].(string)
	if !ok {
		return nil, utils.ErrInvalidVariable
	}

	// 2. 构造MCP请求参数
	argsMap := make(map[string]interface{})
	params := tool["parameters"].([]interface{})
	for i, param := range params {
		paramDef := param.(map[string]interface{})
		paramName := paramDef["name"].(string)
		if i < len(args) {
			argsMap[paramName] = args[i]
		} else if defaultValue, ok := paramDef["default"]; ok {
			argsMap[paramName] = defaultValue
		}
	}

	// 3. 调用MCP工具
	req := map[string]interface{}{
		"server_name": server,
		"tool_name":   toolName,
		"arguments":   argsMap,
	}

	// 这里是简化实现，实际应该通过MCP客户端调用
	res, err := callMCP(req)
	if err != nil {
		return nil, fmt.Errorf("mcp call %s/%s failed: %v", server, toolName, err)
	}

	return res, nil
}

// invokeGRPCMethod 通过反射执行GRPC调用
func invokeGRPCMethod(
	ctx context.Context,
	conn *grpc.ClientConn,
	service, method string,
	reqBytes []byte,
) (interface{}, error) {
	// 简化的反射调用实现
	// 实际项目中需要考虑更完善的反射处理
	md := make(metadata.MD)
	resp := new(bytes.Buffer)

	err := conn.Invoke(
		ctx,
		fmt.Sprintf("/%s/%s", service, method),
		reqBytes,
		resp,
		grpc.Header(&md),
	)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Manager) Render(promptID string, variables map[string]interface{}) (string, error) {
	// 1. 查找模板
	t, ok := m.templates[promptID]
	if !ok {
		return "", utils.ErrTemplateNotFound
	}

	// 2. 验证必需变量
	if err := validateVariables(promptID, variables); err != nil {
		return "", err
	}

	// 3. 合并上下文数据: 变量 + 环境变量
	data := make(map[string]interface{})
	for _, v := range variables {
		data["variables"] = v
	}
	for jsonPath, val := range m.environs {
		data[jsonPath] = val
	}

	// 4. 执行渲染(带超时控制)
	var buf strings.Builder
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 使用channel捕获渲染结果
	resultCh := make(chan string)
	errCh := make(chan error)

	go func() {
		err := t.Execute(&buf, data)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- buf.String()
	}()

	select {
	case <-ctx.Done():
		return "", utils.ErrRenderTimeout
	case err := <-errCh:
		return "", err
	case result := <-resultCh:
		return result, nil
	}
}

// validateVariables 验证模板变量是否符合要求
// validateVariables 验证模板变量是否符合要求
func validateVariables(promptID string, vars map[string]interface{}) error {
	// 1. 获取模板定义
	key := fmt.Sprintf("shenma:templates:%s", promptID)
	// 2. 解析模板定义中的参数要求
	var def map[string]interface{}
	if err := dao.GetJSON(key, &def); err != nil {
		return fmt.Errorf("invalid template definition: %v", err)
	}

	parameters, ok := def["parameters"].([]interface{})
	if !ok {
		// 没有参数定义
		return nil
	}

	// 3. 遍历参数规则进行验证
	for _, param := range parameters {
		paramDef := param.(map[string]interface{})
		paramName := paramDef["name"].(string)

		// 检查必填参数
		if required, ok := paramDef["required"].(bool); ok && required {
			if _, exists := vars[paramName]; !exists {
				return fmt.Errorf("%w: required parameter '%s' is missing",
					utils.ErrInvalidVariable, paramName)
			}
		}

		// 如果有值则验证类型和格式
		if val, exists := vars[paramName]; exists {
			// 类型验证
			if typeStr, ok := paramDef["type"].(string); ok {
				switch typeStr {
				case "string":
					if _, ok := val.(string); !ok {
						return fmt.Errorf("%w: parameter '%s' must be string type",
							utils.ErrInvalidVariable, paramName)
					}
				case "number":
					if _, ok := val.(float64); !ok {
						return fmt.Errorf("%w: parameter '%s' must be number type",
							utils.ErrInvalidVariable, paramName)
					}
				case "boolean":
					if _, ok := val.(bool); !ok {
						return fmt.Errorf("%w: parameter '%s' must be boolean type",
							utils.ErrInvalidVariable, paramName)
					}
				case "array":
					if _, ok := val.([]interface{}); !ok {
						return fmt.Errorf("%w: parameter '%s' must be array type",
							utils.ErrInvalidVariable, paramName)
					}
				case "object":
					if _, ok := val.(map[string]interface{}); !ok {
						return fmt.Errorf("%w: parameter '%s' must be object type",
							utils.ErrInvalidVariable, paramName)
					}
				}
			}

			// 枚举值验证
			if enumValues, ok := paramDef["enum"].([]interface{}); ok {
				found := false
				for _, ev := range enumValues {
					if ev == val {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("%w: parameter '%s' must be one of %v",
						utils.ErrInvalidVariable, paramName, enumValues)
				}
			}

			// 正则表达式验证
			if pattern, ok := paramDef["pattern"].(string); ok {
				matched, err := regexp.MatchString(pattern, val.(string))
				if err != nil {
					return fmt.Errorf("%w: invalid pattern for parameter '%s': %v",
						utils.ErrInvalidVariable, paramName, err)
				}
				if !matched {
					return fmt.Errorf("%w: parameter '%s' must match pattern '%s'",
						utils.ErrInvalidVariable, paramName, pattern)
				}
			}
		}
	}
	return nil
}
