// Package router contains the routes for the API
package router

import (
	queue "jobd/controllers/queue"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	gin.ForceConsoleColor()

	r := gin.Default()
	r.POST("/api/upload", queue.UploadJob)
	r.GET("/api/get/:id", queue.RetrieveJob)

	return r
}
