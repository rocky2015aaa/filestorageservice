package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rocky2015aaa/filestorageservice/internal/api/handlers"
	"github.com/rocky2015aaa/filestorageservice/internal/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const preflightCacheMaxAge = 12 * time.Hour

func NewRouter(handler *handlers.Handler) http.Handler {
	gin.SetMode(os.Getenv(config.EnvSvrGinMode))
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := router.Group("/api/v1")

	v1.GET("/ping", handler.Ping)
	v1.POST("/upload", handler.UploadHandler)
	v1.GET("/files-data", handler.GetFilesDataHandler)
	v1.GET("/download", handler.DownloadHandler)

	router.Use(corsMiddleware())

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "Method Not Allowed",
		})
	})

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
			"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           preflightCacheMaxAge,
	})
}
