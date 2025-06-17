package dao

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ToolCache struct {
	tools map[string]Tool
}

/**
 * Create new ToolCache instance
 * @return Pointer to initialized ToolCache
 */
func NewToolCache() *ToolCache {
	c := &ToolCache{
		tools: make(map[string]Tool),
	}
	return c
}

/**
 * Get all tools from cache
 * @param c ToolCache instance
 * @return Map of all tools
 */
func (c *ToolCache) All() map[string]Tool {
	return c.tools
}

/**
 * Register tool in cache
 * @param c ToolCache instance
 * @param toolId ID of the tool
 * @param tool Tool definition
 */
func (c *ToolCache) Register(toolId string, tool Tool) {
	if toolId == "" {
		logrus.Errorf("Tool ID cannot be empty")
		return
	}
	c.tools[toolId] = tool
}

/**
 * Get tool from cache by ID
 * @param c ToolCache instance
 * @param toolId ID of the tool
 * @return Tool details and exists flag
 */
func (c *ToolCache) Get(toolId string) (Tool, bool) {
	tool, ok := c.tools[toolId]
	return tool, ok
}

/**
 * Load tools from Redis into cache
 * @param c ToolCache instance
 * @param ctx Context for Redis operations
 * @return Error if loading fails
 */
func (c *ToolCache) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading tools from Redis")

	keys, err := KeysByPrefix(PREFIX_TOOLS)
	if err != nil {
		return err
	}

	newTools := make(map[string]Tool)
	for _, key := range keys {
		var tool Tool
		err := GetJSON(key, &tool)
		if err != nil {
			return err
		}
		toolId := KeyToID(key, PREFIX_TOOLS)
		if toolId == "" {
			logrus.Errorf("Tool ID cannot be empty")
			continue
		}
		newTools[toolId] = tool
	}
	c.tools = newTools
	return nil
}
