package service

import (
	"ai-prompt-shell/dao"
	"sync"
	"text/template"
	"time"
)

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

type Renderer struct {
	refreshInterval time.Duration
	templates       map[string]*template.Template
	funcMap         template.FuncMap
	stats           *RenderStats
}

var renderer *Renderer = NewRenderer()

func NewRenderer() *Renderer {
	return &Renderer{
		refreshInterval: 10 * time.Second,
		templates:       make(map[string]*template.Template),
		funcMap:         make(template.FuncMap),
		stats:           &RenderStats{},
	}
}

func RenderPrompt(prompt_id string, variables map[string]interface{}) (string, error) {
	// prompt, _ := prompts.Get(prompt_id)
	// return prompt.Render(variables)

	// 2. 渲染模板
	// prompt, err := tmplManager.Render(promptID, req.Variables)
	// if err != nil {
	// 	if err == utils.ErrTemplateNotFound {
	// 		c.JSON(http.StatusNotFound, gin.H{
	// 			"error": "prompt template not found",
	// 		})
	// 	} else {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"error": "failed to render template",
	// 		})
	// 	}
	// 	return
	// }

	// // 3. 调用LLM (重试2次)
	// var resp llm.ChatCompletionResponse
	// var lastErr error

	// for i := 0; i < 3; i++ {
	// 	if i > 0 {
	// 		time.Sleep(time.Duration(i*100) * time.Millisecond)
	// 	}

	// 	llmReq := llm.ChatCompletionRequest{
	// 		Model: req.Model,
	// 		Messages: []llm.ChatMessage{
	// 			{
	// 				Role:    "user",
	// 				Content: prompt,
	// 			},
	// 		},
	// 	}

	// 	resp, lastErr = llmClient.ChatCompletion(c.Request.Context(), llmReq)
	// 	if lastErr == nil {
	// 		break
	// 	}
	// }

	// if lastErr != nil {
	// 	c.JSON(http.StatusBadGateway, gin.H{
	// 		"error":   "LLM服务不可用",
	// 		"details": lastErr.Error(),
	// 	})
	// 	return
	// }
	return "", nil
}

func onRefreshTools() {
	// for _, t := range registry.All() {
	// 	//TODO:
	// }
}

func onRefreshPrompts() {
	for key, content := range prompts.All() {
		t, err := template.New(key).Funcs(renderer.funcMap).Parse(content.UserPrompt)
		if err != nil {
			continue
		}
		renderer.templates[key] = t
	}
}

func RefreshTool(t *dao.Tool) {

}
