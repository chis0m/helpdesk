package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Timestamp time.Time   `json:"timestamp"`
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Timestamp time.Time   `json:"timestamp"`
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

func Success(ctx *gin.Context, code int, data interface{}, message ...string) {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	obj := SuccessResponse{
		Timestamp: time.Now(),
		Code:      code,
		Message:   msg,
	}
	if data != nil {
		obj.Data = data
	}
	ctx.JSON(code, obj)
}

func Failure(ctx *gin.Context, code int, err interface{}, message ...string) {
	msg := "failed"
	if len(message) > 0 {
		msg = message[0]
	}

	if val, ok := err.(error); ok {
		err = val.Error()
	}

	obj := ErrorResponse{
		Timestamp: time.Now(),
		Code:      code,
		Message:   msg,
		Error:     err,
	}
	ctx.JSON(code, obj)
}

func FailureWithAbort(ctx *gin.Context, code int, err interface{}, message ...string) {
	msg := "failed"
	if len(message) > 0 {
		msg = message[0]
	}

	if val, ok := err.(error); ok {
		err = val.Error()
	}

	obj := ErrorResponse{
		Timestamp: time.Now(),
		Code:      code,
		Message:   msg,
		Error:     err,
	}
	ctx.AbortWithStatusJSON(code, obj)
}
