package api

import (
	"github.com/zgsm-ai/ai-prompt-shell/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListEnvirons Get environment variables list
// @Summary List all environment variables
// @Description Get all defined environment variables in system
// @Tags Environs
// @Produce json
// @Success 200 {array} string
// @Router /api/environs [get]
func ListEnvirons(c *gin.Context) {
	envs, err := service.Environments().Keys()
	if err != nil {
		respErrorf(c, http.StatusInternalServerError, "Failed to load environment variables")
		return
	}
	respOK(c, envs)
}

// GetEnviron Get a single environment variable value
// @Summary Get environment variable
// @Description Get value of specified environment variable
// @Tags Environs
// @Produce json
// @Param environ_id path string true "Environment variable ID"
// @Success 200 {object} interface{}
// @Failure 404 {object} ResponseData
// @Router /api/environs/{environ_id} [get]
func GetEnviron(c *gin.Context) {
	environID := c.Param("environ_id")
	val, ok := service.Environments().Get(environID)
	if !ok {
		respErrorf(c, http.StatusNotFound, "Environment variable not found")
		return
	}
	respOK(c, val)
}
