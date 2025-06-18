package api

import (
	"ai-prompt-shell/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListExtensions list all prompt extension IDs
// @Summary List all prompt extension IDs
// @Description Get all available prompt extension IDs in the system
// @Tags Extensions
// @Produce json
// @Success 200 {array} string
// @Router /api/extensions [get]
func ListExtensions(c *gin.Context) {
	extensions, err := service.ExtensionIDs()
	if err != nil {
		respErrorf(c, http.StatusInternalServerError, "failed to load extensions")
		return
	}
	respOK(c, extensions)
}

// GetExtensionDetail get prompt extension details
// @Summary Get specified prompt extension details
// @Description Get detailed information of prompt extension by ID
// @Tags Extensions
// @Produce json
// @Param extension_id path string true "Extension ID"
// @Success 200 {object} dao.PromptExtension
// @Failure 404 {object} ResponseData
// @Router /api/extensions/{extension_id} [get]
func GetExtensionDetail(c *gin.Context) {
	extensionID := c.Param("extension_id")

	ext, exists := service.Extension(extensionID)

	if !exists {
		respErrorf(c, http.StatusNotFound, "extension not found")
		return
	}
	respOK(c, ext)
}
