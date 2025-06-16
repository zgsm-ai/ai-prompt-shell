package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 设置API路由和Swagger文档路由
func SetupRoutes(r *gin.Engine) {
	// 添加swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API分组
	api := r.Group("/api")
	{
		api.GET("/extensions", ListExtensions)
		api.GET("/extensions/:extension_id", GetExtensionDetail)
		api.GET("/prompts", ListPrompts)
		api.GET("/prompts/:prompt_id", GetPromptDetail)
		api.GET("/tools", ListTools)
		api.GET("/tools/:tool_id", GetToolDetail)
		// 环境变量路由
		api.GET("/environs", ListEnvirons)
		api.GET("/environs/:environ_id", GetEnviron)
		api.POST("/render/prompts/:prompt_id", RenderPrompt)
		api.POST("/chat/prompts/:prompt_id", ChatWithPrompt)
	}
}
