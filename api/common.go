package api

import (
	"github.com/zgsm-ai/ai-prompt-shell/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/**
 * API response structure
 */
type ResponseData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

/**
 * Normal API response
 */
func respOK(c *gin.Context, data any) {
	logrus.Debugf("request: %+v, response: %+v", c.Request.RequestURI, data)
	c.JSON(http.StatusOK, data)
	// c.JSON(http.StatusOK, ResponseData{
	// 	Code:    "0",
	// 	Message: "OK",
	// 	Success: true,
	// 	Data:    data,
	// })
}

/**
 * Error API response
 */
func respError(c *gin.Context, code int, err error) {
	logrus.Errorf("request: %+v, error: %s", c.Request.RequestURI, err.Error())
	if httpErr, ok := err.(*utils.HttpError); ok {
		c.JSON(httpErr.Code(), ResponseData{
			Code:    strconv.Itoa(httpErr.Code()),
			Message: httpErr.Error(),
			Success: false,
		})
	} else {
		c.JSON(code, ResponseData{
			Code:    strconv.Itoa(code),
			Message: err.Error(),
			Success: false,
		})
	}
}

func respErrorf(c *gin.Context, code int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)

	logrus.Errorf("request: %+v, error: %s", c.Request.RequestURI, message)
	c.JSON(code, ResponseData{
		Code:    strconv.Itoa(code),
		Message: message,
		Success: false,
	})
}
