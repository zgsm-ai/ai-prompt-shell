// @title AI Prompt Shell API
// @version 1.0
// @description This is the API documentation for AI Prompt Shell
// @BasePath /
package main

import (
	"github.com/zgsm-ai/ai-prompt-shell/api"
	"github.com/zgsm-ai/ai-prompt-shell/dao"
	"github.com/zgsm-ai/ai-prompt-shell/internal/config"
	"github.com/zgsm-ai/ai-prompt-shell/internal/logger"
	"github.com/zgsm-ai/ai-prompt-shell/service"
	"fmt"

	"log"

	_ "github.com/zgsm-ai/ai-prompt-shell/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	printVersions()

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

var SoftwareVer = ""
var BuildTime = ""
var BuildTag = ""
var BuildCommitId = ""

/*
 * Print software version information
 */
func printVersions() {
	fmt.Printf("Version %s\n", SoftwareVer)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Build Tag: %s\n", BuildTag)
	fmt.Printf("Build Commit ID: %s\n", BuildCommitId)
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
