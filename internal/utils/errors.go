package utils

import (
	"net/http"
)

var (
	ErrKeyNotFound     = NewHttpError(http.StatusNotFound, "key not found")
	ErrPromptNotFound  = NewHttpError(http.StatusNotFound, "prompt not found")
	ErrEnvironNotFound = NewHttpError(http.StatusNotFound, "environment not found")
	ErrToolNotFound    = NewHttpError(http.StatusNotFound, "tool not found")
	ErrRedisError      = NewHttpError(http.StatusInternalServerError, "redis error")
	ErrPromptInvalid   = NewHttpError(http.StatusInternalServerError, "prompt invalid")
	ErrRenderTimeout   = NewHttpError(http.StatusGatewayTimeout, "render timeout")
	ErrToolCallFailed  = NewHttpError(http.StatusInternalServerError, "tool call failed")
	ErrBug             = NewHttpError(http.StatusInternalServerError, "bug")
)

/**
 * Redis operation error wrapper
 */
type RedisError struct {
	Operation string
	Key       string
	Err       error
}

func (e *RedisError) Error() string {
	return "redis error on " + e.Operation + " for key " + e.Key + ": " + e.Err.Error()
}
