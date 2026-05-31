package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Data: data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(errorHTTPStatus(code), Response{Code: code, Message: message})
}

func errorHTTPStatus(code int) int {
	switch code {
	case 40001, 40002, 40003:
		return http.StatusBadRequest
	case 40301:
		return http.StatusForbidden
	case 40401:
		return http.StatusNotFound
	case 40901:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
