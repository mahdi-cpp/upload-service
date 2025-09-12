package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/upload-service/internal/api/handler"
	"github.com/mahdi-cpp/upload-service/internal/application"
)

func main() {

	// Load HTML templates
	Router.LoadHTMLGlob("/app/tmp/templates/*")

	// Create upload handler
	uploadHandler := &handler.UploadHandler{
		UploadDir: "/app/iris/com.iris.settings/uploads",
	}
	// Setup routes
	setupRoutes(Router, uploadHandler)

	newAppManager, err := application.NewAppManager()
	if err != nil {
		log.Fatal(err)
	}

	downloadHandler := handler.NewDownloadHandler(newAppManager)
	routDownloadHandler(downloadHandler)

	startServer(Router)
}

func setupRoutes(router *gin.Engine, uploadHandler *handler.UploadHandler) {
	// Serve upload form
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Setup upload routes
	routUploadHandler(router, uploadHandler)
}

func routUploadHandler(router *gin.Engine, uploadHandler *handler.UploadHandler) {

	router.POST("/api/v1/upload/create", uploadHandler.CreateDirectory)
	router.POST("/api/v1/upload/media", uploadHandler.UploadMedia)
	router.POST("/api/v1/upload/multiple", uploadHandler.UploadMultiple)
	router.GET("/api/v1/upload/files", uploadHandler.ListFiles)
}

func routDownloadHandler(userHandler *handler.DownloadHandler) {

	api := Router.Group("/api/v1/download")

	api.GET("original/*filename", userHandler.ImageOriginal)
	api.GET("thumbnail/*filename", userHandler.ImageThumbnail)
	api.GET("icon/*filename", userHandler.ImageIcons)
}
