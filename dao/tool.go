package dao

// Tool 定义了工具的元数据结构
type Tool struct {
	Name        string                 `json:"name"`
	Module      string                 `json:"module"`
	Type        string                 `json:"type"`
	URL         string                 `json:"url"`
	Description string                 `json:"description"`
	Supports    []string               `json:"supports"`
	Parameters  map[string]interface{} `json:"parameters"`
	Returns     map[string]interface{} `json:"returns"`
	Examples    []string               `json:"examples"`
}

type Restful struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

type Grpc struct {
	Service string `json:"service"`
	Method  string `json:"method"`
}

// ValidToolTypes 定义了有效的工具类型枚举
var ValidToolTypes = []string{"restful", "grpc", "mcp"}

// ValidSupportTypes 定义了有效的支持类型枚举
var ValidSupportTypes = []string{"chat", "completion", "codereview"}
