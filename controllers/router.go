// Package router contains the routes for the API
package router

import (
	queue "jobd/controllers/queue"

	// Import your local docs package
	_ "jobd/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	gin.ForceConsoleColor()

	r := gin.Default()
	r.POST("/api/upload", queue.UploadJob)
	r.GET("/api/get/:id", queue.RetrieveJob)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
