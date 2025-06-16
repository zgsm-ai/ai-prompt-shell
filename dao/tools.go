package dao

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ToolCache struct {
	tools map[string]Tool
}

func NewToolCache() *ToolCache {
	c := &ToolCache{
		tools: make(map[string]Tool),
	}
	return c
}

func (c *ToolCache) All() map[string]Tool {
	return c.tools
}

func (c *ToolCache) Register(toolId string, tool Tool) {
	c.tools[toolId] = tool
}

func (c *ToolCache) Get(toolId string) (Tool, bool) {
	tool, ok := c.tools[toolId]
	return tool, ok
}

func (c *ToolCache) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading tools from Redis")

	keys, err := KeysByPrefix(PREFIX_TOOLS)
	if err != nil {
		return err
	}

	for _, key := range keys {
		var tool Tool
		err := GetJSON(key, &tool)
		if err != nil {
			return err
		}
		c.Register(KeyToPath(key, PREFIX_TOOLS), tool)
	}
	return nil
}
