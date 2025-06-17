package dao

// Tool defines metadata structure for tools
type Tool struct {
	Name        string                 `json:"name"`
	Module      string                 `json:"module"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Supports    []string               `json:"supports"`
	Parameters  map[string]interface{} `json:"parameters"`
	Returns     map[string]interface{} `json:"returns"`
	Examples    []string               `json:"examples,omitempty"`
	Restful     *Restful               `json:"restful,omitempty"`
	Grpc        *Grpc                  `json:"grpc,omitempty"`
}

type Restful struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

type Grpc struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

// ValidToolTypes defines valid tool type enums
var ValidToolTypes = []string{"restful", "grpc", "mcp"}

// ValidSupportTypes defines valid support type enums
var ValidSupportTypes = []string{"chat", "completion", "codereview"}
