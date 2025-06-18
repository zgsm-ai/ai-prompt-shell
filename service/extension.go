package service

import (
	"github.com/zgsm-ai/ai-prompt-shell/dao"
	"context"
)

var extensions = dao.NewExtensionCache()

/**
 * Get extension by ID from cache
 * @param extension_id ID of the extension to retrieve
 * @return extension content if found
 * @return bool indicating if extension exists
 */
func Extension(extension_id string) (dao.PromptExtension, bool) {
	return extensions.Get(extension_id)
}

/**
 * Get all available extension IDs
 * @return slice of extension IDs
 * @return error if failed to load from Redis
 */
func ExtensionIDs() ([]string, error) {
	if err := extensions.LoadFromRedis(context.Background()); err != nil {
		return nil, err
	}
	var result []string
	for k, _ := range extensions.All() {
		result = append(result, k)
	}
	return result, nil
}
