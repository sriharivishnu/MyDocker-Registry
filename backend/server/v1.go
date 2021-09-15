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

	auth := v1.Group("auth")
	auth.POST("/signup", authController.SignUp)
	auth.POST("/login", authController.SignIn)

	users := v1.Group("users")
	repositories := users.Group("/:user_id/repositories")
	images := repositories.Group("/:repo_id/images")

	// /repositories
	v1.GET("/repositories/search", repoController.Search)

	// /users/:id/repositories
	repositories.GET("", repoController.GetForUser)
	repositories.GET("/:repo_id", repoController.GetForUser)

	// /users/:id/repositories/:id/images
	images.GET("", imageTagController.GetImageTagsForRepoName)
	images.GET("/:image_id", imageTagController.GetImage)

	// endpoints that require auth
	repositories.Use(middlewares.AuthMiddleware())
	repositories.POST("", repoController.Create)

	images.Use(middlewares.AuthMiddleware())
	images.POST("", imageTagController.CreateImageTag)
	images.GET("/:image_id/upload_url", imageTagController.GetUploadURL)

}
