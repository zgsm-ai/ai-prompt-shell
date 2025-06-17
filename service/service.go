package service

import (
	"ai-prompt-shell/dao"
	"ai-prompt-shell/internal/config"
	"ai-prompt-shell/internal/utils"
	"context"
	"fmt"
	"time"
)

var llmClient *LLMClient

/**
 * Initialize service with configuration
 * @param c configuration containing API keys and refresh intervals
 * @return error if initialization fails (e.g. redis connection)
 */
func Init(c *config.Config) error {
	if dao.Client == nil {
		return utils.ErrRedisError
	}
	llmClient = NewLLMClient(c.LLM.ApiBase, c.LLM.ApiKey)

	extensions.LoadFromRedis(context.Background())
	tools.LoadFromRedis(context.Background())
	environs.LoadFromRedis(context.Background())
	prompts.LoadFromRedis(context.Background())
	onRefreshExtensions()
	onRefreshTools()
	onRefreshPrompts()

	go startAutoRefreshTools(c.Refresh.Tool)
	go startAutoRefreshPrompts(c.Refresh.Prompt)
	go startAutoRefreshExtensions(c.Refresh.Extension)
	go startAutoRefreshEnvirionments(c.Refresh.Environ)
	return nil
}

/**
 * Start periodic refresh of tools from Redis
 * @param interval duration between refreshes
 */
func startAutoRefreshTools(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			tools.LoadFromRedis(ctx)
			onRefreshTools()
			cancel()
		}
	}
}

/**
 * Start periodic refresh of prompts from Redis
 * @param interval duration between refreshes
 */
func startAutoRefreshPrompts(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			prompts.LoadFromRedis(context.Background())
			onRefreshPrompts()
		}
	}
}

/**
 * Start periodic refresh of extensions from Redis
 * @param interval duration between refreshes
 */
func startAutoRefreshExtensions(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			extensions.LoadFromRedis(context.Background())
			onRefreshExtensions()
		}
	}
}

/**
 * Handle extension refresh by updating contributed prompts
 */
func onRefreshExtensions() {
	for _, ext := range extensions.All() {
		for _, p := range ext.Contributes.Prompts {
			prompt_id := fmt.Sprintf("%s.%s", ext.Name, p.Name)
			_, origin := prompts.Get(prompt_id)
			if origin != dao.PromptOrigin_Direct {
				prompts.Set(prompt_id, p, dao.PromptOrigin_Extension)
			}
		}
	}
}

/**
 * Start periodic refresh of environments from Redis
 * @param interval duration between refreshes
 */
func startAutoRefreshEnvirionments(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			environs.LoadFromRedis(context.Background())
		}
	}
}
