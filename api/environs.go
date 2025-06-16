package api

import (
	"ai-prompt-shell/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListEnvirons 获取环境变量列表
// @Summary 列出所有环境变量
// @Description 获取系统中定义的所有环境变量
// @Tags Environs
// @Produce json
// @Success 200 {object} []string
// @Router /api/environs [get]
func ListEnvirons(c *gin.Context) {
	envs, err := service.Environments().Keys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加载环境变量失败"})
		return
	}
	c.JSON(http.StatusOK, envs)
}

// GetEnviron 获取单个环境变量的值
// @Summary 获取环境变量
// @Description 获取指定环境变量的值
// @Tags Environs
// @Produce json
// @Param environ_id path string true "环境变量ID"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Router /api/environs/{environ_id} [get]
func GetEnviron(c *gin.Context) {
	environID := c.Param("environ_id")
	val, ok := service.Environments().Get(environID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "环境变量不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": val})
}
