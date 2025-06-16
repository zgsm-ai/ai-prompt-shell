package dao

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ExtensionCache struct {
	extensions map[string]PromptExtension
}

func NewExtensionCache() *ExtensionCache {
	c := &ExtensionCache{
		extensions: make(map[string]PromptExtension),
	}
	return c
}

func (c *ExtensionCache) Set(extension_id string, value PromptExtension) {
	c.extensions[extension_id] = value
}

func (c *ExtensionCache) Get(extension_id string) (PromptExtension, bool) {
	val, ok := c.extensions[extension_id]
	return val, ok
}

func (c *ExtensionCache) All() map[string]PromptExtension {
	return c.extensions
}

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
		newExts[extension_id] = val
	}

	c.extensions = newExts
	return nil
}
