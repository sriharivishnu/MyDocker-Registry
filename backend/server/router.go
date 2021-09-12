package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "DELETE", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowWildcard:    true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// health checks
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from shopify Challenge API"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Status OK"})
	})

	SetUpV1(router)
	return router
}
