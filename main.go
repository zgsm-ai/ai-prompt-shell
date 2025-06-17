// @title AI Prompt Shell API
// @version 1.0
// @description This is the API documentation for AI Prompt Shell
// @host localhost:8080
// @BasePath /
package main

import (
	"ai-prompt-shell/api"
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/config"
	"ai-prompt-shell/internal/logger"
	"ai-prompt-shell/service"

	"log"

	_ "ai-prompt-shell/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()
	logger.Init(&cfg.Logger)

	err := dao.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logrus.Fatalf("Redis initialization failed: %v", err)
	}
	if err := service.Init(cfg); err != nil {
		logrus.Fatalf("Service initialization failed: %v", err)
	}
	runHttpServer(&cfg.Server)
}

/*
 * Start HTTP server and register routes
 */
func runHttpServer(c *config.ServerConfig) {
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	api.SetupRoutes(r)

	err := r.Run(c.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
