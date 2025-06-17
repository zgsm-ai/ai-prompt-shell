package dao

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

/**
 * Initialize Redis client connection
 * @param addr Redis server address
 * @param password Redis auth password
 * @param db Redis database number
 * @return Error if connection fails
 */
func InitRedis(addr, password string, db int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if Client == nil {
		return errors.New("failed to create redis client")
	}
	// Test connection
	if _, err := Client.Ping(Ctx).Result(); err != nil {
		Client = nil
		return errors.Wrap(err, "failed to connect to redis")
	}
	return nil
}

/**
 * Set JSON encoded value in Redis
 * @param key Redis key
 * @param value Value to be stored
 * @param expiration Key expiration duration
 * @return Error if operation fails
 */
func SetJSON(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "failed to marshal value")
	}
	return Client.Set(Ctx, key, data, expiration).Err()
}

/**
 * Get JSON decoded value from Redis
 * @param key Redis key
 * @param dest Destination object to store data
 * @return Error if operation fails
 */
func GetJSON(key string, dest any) error {
	data, err := Client.Get(Ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return errors.Wrap(err, "failed to get value")
	}
	return json.Unmarshal(data, dest)
}

/**
 * Delete key from Redis
 * @param key Redis key to delete
 * @return Error if operation fails
 */
func Del(key string) error {
	return Client.Del(Ctx, key).Err()
}

/**
 * Check if key exists in Redis
 * @param key Redis key to check
 * @return exists Whether key exists
 * @return Error if operation fails
 */
func Exists(key string) (bool, error) {
	n, err := Client.Exists(Ctx, key).Result()
	return n > 0, err
}

/**
 * Find keys matching prefix pattern
 * @param prefix Key prefix pattern
 * @return Matching keys list
 * @return Error if operation fails
 */
func KeysByPrefix(prefix string) ([]string, error) {
	var keys []string
	var cursor uint64
	var err error

	for {
		// Safely iterate keys using SCAN command
		var partialKeys []string
		partialKeys, cursor, err = Client.Scan(Ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan keys")
		}

		keys = append(keys, partialKeys...)

		if cursor == 0 { // Iteration completed
			break
		}
	}

	return keys, nil
}

/**
 * Load all JSON values under prefix from Redis
 * @param prefix Key prefix pattern
 * @return Map of key-value pairs
 * @return Error if operation fails
 */
func LoadJsons(prefix string) (map[string]interface{}, error) {
	keys, err := KeysByPrefix(prefix)
	if err != nil {
		return nil, err
	}

	jsons := make(map[string]interface{})
	for _, key := range keys {
		var val interface{}
		if err := GetJSON(key, &val); err != nil {
			return nil, err
		}
		jsons[key] = val
	}

	return jsons, nil
}
