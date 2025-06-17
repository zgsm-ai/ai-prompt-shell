package dao

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ExtensionCache struct {
	extensions map[string]PromptExtension
}

/**
 * Create new ExtensionCache instance
 * @return Pointer to initialized ExtensionCache
 */
func NewExtensionCache() *ExtensionCache {
	c := &ExtensionCache{
		extensions: make(map[string]PromptExtension),
	}
	return c
}

/**
 * Set extension in cache
 * @param c ExtensionCache instance
 * @param extension_id Extension ID
 * @param value Extension details
 */
func (c *ExtensionCache) Set(extension_id string, value PromptExtension) {
	c.extensions[extension_id] = value
}

/**
 * Get extension from cache
 * @param c ExtensionCache instance
 * @param extension_id Extension ID
 * @return Extension details and exists flag
 */
func (c *ExtensionCache) Get(extension_id string) (PromptExtension, bool) {
	val, ok := c.extensions[extension_id]
	return val, ok
}

/**
 * Get all extensions from cache
 * @param c ExtensionCache instance
 * @return Map of all extensions
 */
func (c *ExtensionCache) All() map[string]PromptExtension {
	return c.extensions
}

/**
 * Load extensions from Redis into cache
 * @param c ExtensionCache instance
 * @param ctx Context for Redis operations
 * @return Error if loading fails
 */
func (c *ExtensionCache) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading extensions from Redis")

	keys, err := KeysByPrefix(PREFIX_EXTENSIONS)
	if err != nil {
		return err
	}

	newExts := make(map[string]PromptExtension)
	for _, key := range keys {
		var val PromptExtension
		if err := GetJSON(key, &val); err != nil {
			return err
		}

		extension_id := KeyToID(key, PREFIX_EXTENSIONS)
		if extension_id == "" {
			logrus.Warnf("Extension ID cannot be empty")
			continue
		}
		newExts[extension_id] = val
	}

	c.extensions = newExts
	return nil
}
