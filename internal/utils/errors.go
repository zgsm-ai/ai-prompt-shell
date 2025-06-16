package utils

import "errors"

var (
	ErrRedisError       = errors.New("redis error")
	ErrKeyNotFound      = errors.New("key not found")
	ErrTemplateNotFound = errors.New("template not found")
	ErrInvalidVariable  = errors.New("invalid variable format")
	ErrRenderTimeout    = errors.New("render timeout")
	ErrToolCallFailed   = errors.New("tool call failed")
)

type RedisError struct {
	Operation string
	Key       string
	Err       error
}

func (e *RedisError) Error() string {
	return "redis error on " + e.Operation + " for key " + e.Key + ": " + e.Err.Error()
}
