package api

import (
	"ai-prompt-shell/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListExtensions 列出所有Prompt类型扩展的ID
// @Summary 获取所有Prompt类型扩展的ID
// @Description 获取系统中可用的所有Prompt类型扩展的ID
// @Tags Extensions
// @Produce json
// @Success 200 {array} string
// @Router /api/extensions [get]
func ListExtensions(c *gin.Context) {
	extensions, err := service.ExtensionIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load extensions",
		})
		return
	}

	c.JSON(http.StatusOK, extensions)
}

// GetExtensionDetail 获取Prompt类型扩展详情
// @Summary 获取指定Prompt类型扩展详情
// @Description 根据ID获取Prompt类型扩展的详细信息
// @Tags Extensions
// @Produce json
// @Param extension_id path string true "扩展ID"
// @Success 200 {object} dao.PromptExtension
// @Failure 404 {object} map[string]interface{}
// @Router /api/extensions/{extension_id} [get]
func GetExtensionDetail(c *gin.Context) {
	extensionID := c.Param("extension_id")

	ext, exists := service.Extension(extensionID)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "extension not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"extension": ext,
	})
}
