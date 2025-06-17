package dao

import (
	"context"

	"github.com/sirupsen/logrus"
)

type PromptOrigin string

const (
	PromptOrigin_Notexist  PromptOrigin = "notexist"
	PromptOrigin_Direct    PromptOrigin = "direct"
	PromptOrigin_Extension PromptOrigin = "extension"
)

type PromptLoaded struct {
	Prompt
	Origin PromptOrigin
}

/**
 * Cache for storing prompt templates with origins
 */
type PromptCache struct {
	templates map[string]PromptLoaded
}

/**
 * Create new PromptCache instance
 * @return Pointer to initialized PromptCache
 */
func NewPromptCache() *PromptCache {
	c := &PromptCache{
		templates: make(map[string]PromptLoaded),
	}
	return c
}

/**
 * Add or update prompt template in cache
 * @param c PromptCache instance
 * @param prompt_id ID of the prompt template
 * @param value Prompt template content
 * @param origin Source of the prompt (Direct or Extension)
 */
func (c *PromptCache) Set(prompt_id string, value Prompt, origin PromptOrigin) {
	if origin != PromptOrigin_Direct && origin != PromptOrigin_Extension {
		panic("Invalid origin")
	}
	c.templates[prompt_id] = PromptLoaded{
		Prompt: value,
		Origin: origin,
	}
}

/**
 * Get prompt template from cache
 * @param c PromptCache instance
 * @param prompt_id ID of the prompt template
 * @return Prompt template and its origin
 */
func (c *PromptCache) Get(prompt_id string) (Prompt, PromptOrigin) {
	val, ok := c.templates[prompt_id]
	if !ok {
		return Prompt{}, PromptOrigin_Notexist
	}
	return val.Prompt, val.Origin
}

/**
 * Get all prompt templates from cache
 * @param c PromptCache instance
 * @return Map of all prompt templates
 */
func (c *PromptCache) All() map[string]PromptLoaded {
	return c.templates
}

/**
 * Load prompt templates from Redis into cache
 * @param c PromptCache instance
 * @param ctx Context for Redis operations
 * @return Error if loading fails
 */
func (c *PromptCache) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading templates from Redis")

	keys, err := KeysByPrefix(PREFIX_TEMPLATES)
	if err != nil {
		return err
	}

	newPrompts := make(map[string]PromptLoaded)
	//	Migrate extension-registered prompt templates to newPrompts
	for k, t := range c.templates {
		if t.Origin == PromptOrigin_Extension {
			newPrompts[k] = t
		}
	}
	//	Load directly-registered prompt templates from Redis, which may override extension-registered ones
	for _, key := range keys {
		var val Prompt
		if err := GetJSON(key, &val); err != nil {
			return err
		}

		prompt_id := KeyToID(key, PREFIX_TEMPLATES)
		if prompt_id == "" {
			logrus.Warnf("Prompt ID cannot be empty")
			continue
		}
		newPrompts[prompt_id] = PromptLoaded{
			Prompt: val,
			Origin: PromptOrigin_Direct,
		}
	}

	c.templates = newPrompts
	return nil
}
