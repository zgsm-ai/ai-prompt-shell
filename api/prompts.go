package api

import (
	"net/http"

	"ai-prompt-shell/dao"
	"ai-prompt-shell/service"

	"github.com/gin-gonic/gin"
)

// ListPrompts 列出所有Prompt模板
// @Summary 获取所有Prompt模板
// @Description 获取系统中可用的所有Prompt模板列表
// @Tags Prompts
// @Produce json
// @Success 200 {array} string
// @Router /api/prompts [get]
func ListPrompts(c *gin.Context) {
	templates, err := service.PromptIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load templates",
		})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetPromptDetail 获取Prompt模板详情
// @Summary 获取指定Prompt模板详情
// @Description 根据ID获取Prompt模板的详细信息
// @Tags Prompts
// @Produce json
// @Param prompt_id path string true "Prompt模板ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/prompts/{prompt_id} [get]
func GetPromptDetail(c *gin.Context) {
	promptID := c.Param("prompt_id")

	if dao.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to connect redis",
		})
		return
	}
	prompt, origin := service.Prompt(promptID)
	if origin == dao.PromptOrigin_Notexist {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "template not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"origin": string(origin),
		"prompt": prompt,
	})
}

// RenderPrompt 渲染Prompt模板
// @Summary 渲染指定Prompt模板
// @Description 使用给定变量渲染指定的Prompt模板
// @Tags Render
// @Accept json
// @Produce json
// @Param prompt_id path string true "Prompt模板ID"
// @Param request body map[string]interface{} true "渲染参数"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/render/prompts/{prompt_id} [post]
func RenderPrompt(c *gin.Context) {
	promptID := c.Param("prompt_id")
	var req struct {
		Variables map[string]interface{} `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	kind, data, err := service.RenderPrompt(promptID, req.Variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to render template",
		})
		return
	}
	if kind == "prompt" {
		c.JSON(http.StatusOK, gin.H{
			"kind":   kind,
			"prompt": data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"kind":     kind,
			"messages": data,
		})
	}
}

type ChatModelRequest struct {
	Model     string                 `json:"model"`
	Variables map[string]interface{} `json:"variables"`
}

// ChatWithPrompt 使用Prompt与LLM聊天
// @Summary 使用Prompt与LLM交互
// @Description 使用指定的Prompt模板与LLM进行聊天交互
// @Tags Chat
// @Accept json
// @Produce json
// @Param prompt_id path string true "Prompt模板ID"
// @Param request body ChatModelRequest true "聊天参数"
// @Success 200 {object} service.ChatResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/chat/prompts/{prompt_id} [post]
func ChatWithPrompt(c *gin.Context) {
	promptID := c.Param("prompt_id")

	var req ChatModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	resp, err := service.ChatWithPrompt(promptID, req.Model, req.Variables)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
