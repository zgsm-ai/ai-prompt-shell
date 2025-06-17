package api

import (
	"net/http"

	"ai-prompt-shell/dao"
	"ai-prompt-shell/service"

	"github.com/gin-gonic/gin"
)

// ListPrompts list all prompt templates
// @Summary List all prompt templates
// @Description Get all available prompt templates in the system
// @Tags Prompts
// @Produce json
// @Success 200 {array} string
// @Router /api/prompts [get]
func ListPrompts(c *gin.Context) {
	templates, err := service.PromptIDs()
	if err != nil {
		respErrorf(c, http.StatusInternalServerError, "failed to load templates")
		return
	}
	respOK(c, templates)
}

// GetPromptDetail get prompt template details
// @Summary Get specified prompt template details
// @Description Get detailed information of prompt template by ID
// @Tags Prompts
// @Produce json
// @Param prompt_id path string true "Prompt template ID"
// @Success 200 {object} dao.Prompt
// @Failure 404 {object} map[string]interface{}
// @Router /api/prompts/{prompt_id} [get]
func GetPromptDetail(c *gin.Context) {
	promptID := c.Param("prompt_id")

	if dao.Client == nil {
		respErrorf(c, http.StatusInternalServerError, "failed to connect redis")
		return
	}
	prompt, origin := service.Prompt(promptID)
	if origin == dao.PromptOrigin_Notexist {
		respErrorf(c, http.StatusNotFound, "template not found")
		return
	}

	respOK(c, gin.H{
		"origin": string(origin),
		"prompt": prompt,
	})
}

// RenderPrompt render prompt template
// @Summary Render specified prompt template
// @Description Render the prompt template with given args
// @Tags Prompts
// @Accept json
// @Produce json
// @Param prompt_id path string true "Prompt template ID"
// @Param request body map[string]interface{} true "Rendering parameters"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/prompts/{prompt_id}/render [post]
func RenderPrompt(c *gin.Context) {
	promptID := c.Param("prompt_id")
	var req struct {
		Args map[string]interface{} `json:"args"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respErrorf(c, http.StatusBadRequest, "invalid request body")
		return
	}

	kind, data, err := service.RenderPrompt(promptID, req.Args)
	if err != nil {
		respErrorf(c, http.StatusInternalServerError, "failed to render template")
		return
	}
	if kind == "prompt" {
		respOK(c, gin.H{
			"kind":   kind,
			"prompt": data,
		})
	} else {
		respOK(c, gin.H{
			"kind":     kind,
			"messages": data,
		})
	}
}

type ChatModelRequest struct {
	Model string                 `json:"model"`
	Args  map[string]interface{} `json:"args"`
}

// ChatWithPrompt chat with LLM using prompt
// @Summary Interact with LLM using prompt
// @Description Chat interaction with LLM using specified prompt template
// @Tags Prompts
// @Accept json
// @Produce json
// @Param prompt_id path string true "Prompt template ID"
// @Param request body ChatModelRequest true "Chat parameters"
// @Success 200 {object} service.ChatResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/prompts/{prompt_id}/chat [post]
func ChatWithPrompt(c *gin.Context) {
	promptID := c.Param("prompt_id")

	var req ChatModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respErrorf(c, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := service.ChatWithPrompt(promptID, req.Model, req.Args)
	if err != nil {
		respError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
