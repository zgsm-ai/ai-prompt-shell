package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures API routes and Swagger documentation routes
func SetupRoutes(r *gin.Engine) {
	// Add swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// API group
	api := r.Group("/api")
	{
		api.GET("/extensions", ListExtensions)
		api.GET("/extensions/:extension_id", GetExtensionDetail)
		api.GET("/prompts", ListPrompts)
		api.GET("/prompts/:prompt_id", GetPromptDetail)
		api.POST("/prompts/:prompt_id/render", RenderPrompt)
		api.POST("/prompts/:prompt_id/chat", ChatWithPrompt)
		api.GET("/tools", ListTools)
		api.GET("/tools/:tool_id", GetToolDetail)
		// Environment variables routes
		api.GET("/environs", ListEnvirons)
		api.GET("/environs/:environ_id", GetEnviron)
	}
}
