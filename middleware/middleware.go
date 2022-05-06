package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/* func Cors() gin.HandlerFunc { */
/*     return func(c *gin.Context) { */
/*         c.Writer.Header().Set("Access-Control-Allow-Origin", "*") */
/*         c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin") */
/*         c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT") */
/*         if c.Request.Method == "OPTIONS" { */
/*             c.AbortWithStatus(204) */
/*             return */
/*         } */
/*         defer func() { */
/*             if err := recover(); err != nil { */
/*                 // core.Logger.Error("Panic info is: %v", err) */
/*                 // core.Logger.Error("Panic info is: %s", debug.Stack()) */
/*             } */
/*         }() */
/*         c.Next() */
/*     } */
/* } */

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

/*
func MyTimeMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Println("mytime start.")
		c.Set("request", "中间件")
		c.Next()

		status := c.Writer.Status()
		log.Println("done.", status)

		latency := time.Since(start)
		log.Println("time: ", latency)

		// log.Println(status)
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("lg_cookie"); err == nil {
			if cookie == "good" {
				c.Next()
				return
			}
		}
		// 返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		// 若验证不通过，不再调用后续的函数处理
		c.Abort()
		return
	}
}
*/
