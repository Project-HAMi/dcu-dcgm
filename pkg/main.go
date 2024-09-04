package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "g.sugon.com/das/dcgm-dcu/pkg/docs"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample example
//	@host			127.0.0.1:8080
//	@BasePath		/router/v1

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Group("/")

	r.Run(":16081")
}
