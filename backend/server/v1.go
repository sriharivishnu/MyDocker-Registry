package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/controllers"
	"github.com/sriharivishnu/shopify-challenge/layers"
	"github.com/sriharivishnu/shopify-challenge/middlewares"
	"github.com/sriharivishnu/shopify-challenge/services"
)

func SetUpV1(router *gin.Engine) {
	// set up controllers and inject dependencies
	authController := &controllers.AuthController{
		UserService: &layers.UserService{},
	}
	repoController := &controllers.RepositoryController{
		RepositoryService: &layers.RepositoryService{},
	}
	imageTagController := &controllers.ImageController{
		RepositoryService: &layers.RepositoryService{},
		ImageService:      &layers.ImageService{},
		StorageService:    &services.S3{},
	}

	// Set up routes
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
	images.GET("/:image_id", imageTagController.PullImage)

	// endpoints that require auth
	repositories.Use(middlewares.AuthMiddleware(&layers.UserService{}))
	repositories.POST("", repoController.Create)

	images.Use(middlewares.AuthMiddleware(&layers.UserService{}))
	images.POST("", imageTagController.PushImage)
}
