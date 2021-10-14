package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/controllers"
	"github.com/sriharivishnu/shopify-challenge/external"
	"github.com/sriharivishnu/shopify-challenge/middlewares"
	"github.com/sriharivishnu/shopify-challenge/services"
)

func SetUpV1(router *gin.Engine) {
	// set up controllers and inject dependencies
	authController := &controllers.AuthController{
		UserService: &services.UserService{},
	}
	repoController := &controllers.RepositoryController{
		RepositoryService: &services.RepositoryService{},
	}
	imageTagController := &controllers.ImageController{
		RepositoryService: &services.RepositoryService{},
		ImageService:      &services.ImageService{},
		StorageService:    &external.S3{},
	}

	// Set up routes
	v1 := router.Group("v1")

	auth := v1.Group("auth")
	auth.POST("/signup", authController.SignUp)
	auth.POST("/login", authController.SignIn)

	users := v1.Group("users")
	repositories := users.Group("/:username/repositories")
	images := repositories.Group("/:repo_name/images")

	// /repositories
	v1.GET("/repositories/search", repoController.Search)

	// /users/:id/repositories
	repositories.GET("", repoController.GetForUser)
	// repositories.GET("/:repo_id", repoController.GetForUser)

	// /users/:id/repositories/:id/images
	images.GET("", imageTagController.GetImageTagsForRepoName)
	images.GET("/:image_tag", imageTagController.PullImage)

	// endpoints that require auth
	repositories.Use(middlewares.AuthMiddleware(&services.UserService{}))
	repositories.POST("", repoController.Create)

	images.Use(middlewares.AuthMiddleware(&services.UserService{}))
	images.POST("", imageTagController.PushImage)
}
