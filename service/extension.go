package service

import (
	"ai-prompt-shell/dao"
	"context"
)

var cache = dao.NewExtensionCache()

func Extensions() (map[string]dao.PromptExtension, error) {
	if err := cache.LoadFromRedis(context.Background()); err != nil {
		return nil, err
	}
	return cache.All(), nil
}

func Extension(extension_id string) (dao.PromptExtension, bool) {
	return cache.Get(extension_id)
}

func ExtensionIDs() ([]string, error) {
	if err := cache.LoadFromRedis(context.Background()); err != nil {
		return nil, err
	}
	var result []string
	for k, _ := range cache.All() {
		result = append(result, k)
	}
	return result, nil
}
