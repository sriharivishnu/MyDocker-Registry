package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/controllers"
	"github.com/sriharivishnu/shopify-challenge/middlewares"
)

func SetUpV1(router *gin.Engine) {
	authController := &controllers.AuthController{}
	repoController := &controllers.RepositoryController{}
	imageTagController := &controllers.ImageController{}

	v1 := router.Group("v1")

	v1.POST("/signup", authController.SignUp)
	v1.POST("/signin", authController.SignIn)

	v1.Use(middlewares.AuthMiddleware())

	repository := v1.Group("repositories")
	repository.POST("", repoController.Create)
	repository.GET("", repoController.GetForUser)

	image_tag := v1.Group("tag")
	image_tag.POST("/upload_url", imageTagController.GetUploadURL)

}
