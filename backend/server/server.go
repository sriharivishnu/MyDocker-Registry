package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/config"
)

func Init() {
	config := config.Config
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := NewRouter()
	r.Run(":" + config.PORT)
}
