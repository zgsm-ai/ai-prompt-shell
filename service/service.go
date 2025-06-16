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

func Init(c *config.Config) error {
	if dao.Client == nil {
		return utils.ErrRedisError
	}
	llmClient = NewLLMClient(c.LLM.ApiBase, c.LLM.ApiKey)

	cache.LoadFromRedis(context.Background())
	registry.LoadFromRedis(context.Background())
	env.LoadFromRedis(context.Background())
	prompts.LoadFromRedis(context.Background())

	startAutoRefreshTools(c.Refresh.Tool)
	startAutoRefreshPrompts(c.Refresh.Prompt)
	startAutoRefreshExtensions(c.Refresh.Extension)
	startAutoRefreshEnvirionments(c.Refresh.Environ)
	return nil
}

func startAutoRefreshTools(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			registry.LoadFromRedis(ctx)
			onRefreshTools()
			cancel()
			// case <-c.refreshDone:
			// 	return
		}
	}
}

func startAutoRefreshPrompts(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			prompts.LoadFromRedis(context.Background())
			onRefreshPrompts()
			// case <-c.refreshChannel:
			// Receives manual refresh signal
		}
	}
}

func startAutoRefreshExtensions(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.LoadFromRedis(context.Background())
			onRefreshExtensions()
			// case <-c.refreshChannel:
			// c.onRefreshAll()
			// 收到手动刷新信号
		}
	}
}

func onRefreshExtensions() {
	for _, ext := range cache.All() {
		for _, p := range ext.Contributes.Prompts {
			prompt_id := fmt.Sprintf("%s.%s", ext.Name, p.Name)
			_, origin := prompts.Get(prompt_id)
			if origin != dao.PromptOrigin_Direct {
				prompts.Set(prompt_id, p, dao.PromptOrigin_Extension)
			}
		}
	}
}

func startAutoRefreshEnvirionments(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			env.LoadFromRedis(context.Background())
			// case <-c.refreshChannel:
			// Receives manual refresh signal
		}
	}
}
