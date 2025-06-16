package dao

// PromptExtension represents a prompt extension definition
type PromptExtension struct {
	Name          string      `json:"name" description:"扩展名称"`
	Publisher     string      `json:"publisher" description:"发布者名称"`
	DisplayName   string      `json:"displayName" description:"扩展的显示名称"`
	Icon          string      `json:"icon" description:"扩展的图标路径"`
	Description   string      `json:"description" description:"扩展的描述信息"`
	Version       string      `json:"version" description:"扩展的版本号"`
	ExtensionType string      `json:"extensionType" description:"扩展类型"`
	License       string      `json:"license" description:"扩展的许可协议"`
	Engines       Engines     `json:"engines" description:"引擎配置"`
	Contributes   Contributes `json:"contributes" description:"扩展功能"`
}

// Engines defines the engine requirements
type Engines struct {
	Name    string `json:"name" description:"引擎名称"`
	Version string `json:"version" description:"引擎版本号"`
}

// Contributes defines the extension contributions
type Contributes struct {
	Prompts     []Prompt     `json:"prompts" description:"Prompt模板"`
	Languages   []string     `json:"languages" description:"支持的语言"`
	Dependences []Dependence `json:"dependences" description:"扩展依赖"`
}

// Prompt defines a single prompt template
type Prompt struct {
	Name        string                 `json:"name" description:"Prompt模板名称"`
	Description string                 `json:"description" description:"描述信息"`
	Messages    []Message              `json:"messages,omitempty" description:"消息列表"`
	UserPrompt  string                 `json:"userPrompt,omitempty" description:"用户提示词模板"`
	Supports    []string               `json:"supports" description:"支持的场景"`
	Parameters  map[string]interface{} `json:"parameters" description:"参数定义(JSON Schema)"`
	Returns     map[string]interface{} `json:"returns" description:"返回值定义(JSON Schema)"`
}

// Message defines a role-message pair
type Message struct {
	Role    string `json:"role" description:"消息角色"`
	Content string `json:"content" description:"消息内容"`
}

// Dependence defines an extension dependency
type Dependence struct {
	Name         string `json:"name" description:"依赖名称"`
	Version      string `json:"version" description:"依赖版本号"`
	FailStrategy string `json:"failStrategy" description:"失败策略"`
}

// Constants for enum values
const (
	ExtensionTypePrompt = "prompt"

	MessageRoleSystem = "system"
	MessageRoleUser   = "user"

	SupportChat       = "chat"
	SupportCodeReview = "codereview"

	FailStrategyAbort  = "abort"
	FailStrategyIgnore = "ignore"
)
