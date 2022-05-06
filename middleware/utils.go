package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code string, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}

func Err(c *gin.Context, code string, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}
