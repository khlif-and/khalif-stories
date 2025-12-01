package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, APIResponse{
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Error: message,
	})
}