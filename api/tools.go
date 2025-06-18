package api

import (
	"github.com/zgsm-ai/ai-prompt-shell/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListTools list all available tools
// @Summary List all tools
// @Description Get a list of all available tools in the system
// @Tags Tools
// @Produce json
// @Success 200 {array} string
// @Router /api/tools [get]
func ListTools(c *gin.Context) {
	tools, err := service.ToolIDs()
	if err != nil {
		respError(c, http.StatusInternalServerError, err)
		return
	}
	respOK(c, tools)
}

// GetToolDetail get tool details
// @Summary Get tool details
// @Description Get detailed information about specified tool
// @Tags Tools
// @Produce json
// @Param tool_id path string true "工具ID"
// @Success 200 {object} dao.Tool
// @Failure 404 {object} ResponseData
// @Router /api/tools/{tool_id} [get]
func GetToolDetail(c *gin.Context) {
	toolID := c.Param("tool_id")

	toolDetail, err := service.GetTool(toolID)
	if err != nil {
		respError(c, http.StatusInternalServerError, err)
		return
	}

	respOK(c, toolDetail)
}
