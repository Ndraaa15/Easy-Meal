package helper

import (
	"github.com/gin-gonic/gin"
)

type HTTPResponse struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, code int64, message string, data interface{}) {
	c.JSON(int(code), HTTPResponse{
		Message: message,
		Status:  true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int64, message string, data interface{}) {
	c.JSON(int(code), HTTPResponse{
		Message: message,
		Status:  false,
		Data:    data,
	})
}
