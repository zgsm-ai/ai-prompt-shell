package api

import (
	"ai-prompt-shell/service"

	"github.com/gin-gonic/gin"
)

// ListTools 列出所有可用工具
// @Summary 列出所有工具
// @Description 获取系统中所有可用工具列表
// @Tags Tools
// @Produce json
// @Success 200 {array} string
// @Router /api/tools [get]
func ListTools(c *gin.Context) {
	tools, err := service.ToolIDs()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "获取工具列表失败",
		})
		return
	}
	c.JSON(200, tools)
}

// GetToolDetail 获取工具详情
// @Summary 获取工具详情
// @Description 获取指定工具的详细信息
// @Tags Tools
// @Produce json
// @Param tool_id path string true "工具ID"
// @Success 200 {object} dao.Tool
// @Failure 404 {object} map[string]string
// @Router /api/tools/{tool_id} [get]
func GetToolDetail(c *gin.Context) {
	toolID := c.Param("tool_id")

	toolDetail, err := service.GetTool(toolID)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "工具不存在",
		})
		return
	}

	c.JSON(200, toolDetail)
}
