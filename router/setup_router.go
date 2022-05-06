package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Register func(*gin.Engine)

func Init(routers ...Register) *gin.Engine {
	rs := append([]Register{}, routers...)

	r := gin.Default()
	for _, register := range rs {
		register(r)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
