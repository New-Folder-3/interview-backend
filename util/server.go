package util

import (
	"github.com/gin-gonic/gin"
)

func ErrorResp(c *gin.Context, err string, code int) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

func SuccessResp(c *gin.Context, data interface{}, message string) {
	c.JSON(200, gin.H{
		"message": message,
		"data":    data,
	})
}
