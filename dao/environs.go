package dao

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

type Environments struct {
	environments map[string]interface{}
}

/**
 * Create new Environments instance
 * @return Pointer to initialized Environments
 */
func NewEnvironments() *Environments {
	c := &Environments{
		environments: make(map[string]interface{}),
	}
	return c
}

/**
 * Set environment variable value
 * @param c Environments instance
 * @param key Environment variable key
 * @param value Variable value
 */
func (c *Environments) Set(key string, value interface{}) {
	c.environments[key] = value
}

/**
 * Get environment variable value
 * @param c Environments instance
 * @param key Environment variable key
 * @return Value and exists flag
 */
func (c *Environments) Get(key string) (interface{}, bool) {
	val, ok := c.environments[key]
	return val, ok
}

/**
 * Get all environment variables
 * @param c Environments instance
 * @return Map of all environments
 */
func (c *Environments) All() map[string]interface{} {
	return c.environments
}

/**
 * Get all environment variable keys
 * @param c Environments instance
 * @return Array of keys and error
 */
func (c *Environments) Keys() ([]string, error) {
	var result []string
	for key, _ := range c.environments {
		result = append(result, key)
	}
	return result, nil
}

/**
 * Load environments from Redis
 * @param c Environments instance
 * @param ctx Context for Redis operations
 * @return Error if loading fails
 */
func (c *Environments) LoadFromRedis(ctx context.Context) error {
	logrus.Info("Loading environments from Redis")

	keys, err := KeysByPrefix(PREFIX_ENVIRONS)
	if err != nil {
		return err
	}

	newEnvs := make(map[string]interface{})
	for _, key := range keys {
		var val interface{}
		if err := GetJSON(key, &val); err != nil {
			return err
		}

		jsonPath := KeyToPath(key, PREFIX_ENVIRONS)
		if jsonPath == "" {
			logrus.Warnf("Environ ID cannot be empty")
			continue
		}
		c.setValueByPath(newEnvs, jsonPath, val)
		newEnvs[jsonPath] = val
	}

	c.environments = newEnvs
	return nil
}

/**
 * Set nested value in map by dot-path
 * @param c Environments instance
 * @param data Target map
 * @param path Dot path to target location
 * @param value Value to set
 */
func (c *Environments) setValueByPath(data map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
		} else {
			child, ok := current[part]
			if !ok {
				current[part] = make(map[string]interface{})
			} else {
				if _, ok := child.(map[string]interface{}); !ok {
					current[part] = make(map[string]interface{})
				}
			}
			current = current[part].(map[string]interface{})
		}
	}
}

// func (c *Environments) GetChild(jsonPath string) (interface{}, error) {
// 	parts := strings.Split(jsonPath, ".")
// 	current := c.environments

// 	for i, part := range parts {
// 		if i == len(parts)-1 {
// 			child, ok := current[part]
// 			if !ok {
// 				return child, os.ErrNotExist
// 			}
// 			return child, nil
// 		} else {
// 			child, ok := current[part]
// 			if !ok {
// 				return nil, os.ErrNotExist
// 			}
// 			current, ok = child.(map[string]interface{})
// 			if !ok {
// 				return nil, os.ErrInvalid
// 			}
// 		}
// 	}
// 	return nil, os.ErrNotExist
// }
