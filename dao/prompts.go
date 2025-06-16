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

type PromptCache struct {
	templates map[string]Prompt
	origins   map[string]PromptOrigin
}

func NewPromptCache() *PromptCache {
	c := &PromptCache{
		templates: make(map[string]Prompt),
		origins:   make(map[string]PromptOrigin),
	}
	return c
}

func (c *PromptCache) Set(prompt_id string, value Prompt, origin PromptOrigin) {
	if origin != PromptOrigin_Direct && origin != PromptOrigin_Extension {
		panic("Invalid origin")
	}
	c.templates[prompt_id] = value
	c.origins[prompt_id] = origin
}

func (c *PromptCache) Get(prompt_id string) (Prompt, PromptOrigin) {
	val, ok := c.templates[prompt_id]
	if !ok {
		return val, PromptOrigin_Notexist
	}
	return val, c.origins[prompt_id]
}

func (c *PromptCache) All() map[string]Prompt {
	return c.templates
}

func (c *PromptCache) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading templates from Redis")

	keys, err := KeysByPrefix(PREFIX_TEMPLATES)
	if err != nil {
		return err
	}

	newPrompts := make(map[string]Prompt)
	newOrigins := make(map[string]PromptOrigin)
	for _, key := range keys {
		var val Prompt
		if err := GetJSON(key, &val); err != nil {
			return err
		}

		prompt_id := KeyToID(key, PREFIX_TEMPLATES)
		newPrompts[prompt_id] = val
		newOrigins[prompt_id] = PromptOrigin_Direct
	}
	for k, t := range c.templates {
		origin := c.origins[k]
		if origin == PromptOrigin_Extension {
			newPrompts[k] = t
			newOrigins[k] = origin
		}
	}

	c.templates = newPrompts
	c.origins = newOrigins
	return nil
}
