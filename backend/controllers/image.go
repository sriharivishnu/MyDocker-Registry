package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/layers"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/services"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type ImageController struct {
	RepositoryService layers.RepositoryLayer
	ImageService      layers.ImageLayer
	StorageService    services.Storage
}

func (i *ImageController) PushImage(c *gin.Context) {
	// Input validation
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")
	var payload struct {
		ImageTag    string `json:"tag"`
		Description string `json:"description,omitempty"`
	}
	errInputFormat := c.BindJSON(&payload)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	// Fetch repository to push to
	repo, errGetRepo := i.RepositoryService.GetRepositoryByName(username, repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Could not find repository: "+username+"/"+repoName, http.StatusNotFound)
		return
	}
	if user.Id != repo.OwnerId {
		utils.RespondErrorString(c, "User is not allowed to upload for this repository", http.StatusForbidden)
		return
	}

	// create file key and image
	key := utils.CreateFileKey(user.Username, repo.Name, payload.ImageTag)
	imageTag, errCreate := i.ImageService.Create(repo.Id, payload.ImageTag, payload.Description, key)
	if errCreate != nil {
		utils.RespondSQLError(c, errCreate)
		return
	}

	// get the upload URL for user to push to
	// Ideally, we have another service (lambda) to notify
	// server when upload is successful.
	URL, err := i.StorageService.GetUploadURL(user.Username, repoName, payload.ImageTag)
	if err != nil {
		utils.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Created image successfully", "id": imageTag.Id, "upload_url": URL})
}

func (i *ImageController) PullImage(c *gin.Context) {
	// input parameters
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")
	imageName, _ := c.Params.Get("image_id")
	if imageName == "" {
		imageName = "latest"
	}

	// Fetch repository
	repo, errGetRep := i.RepositoryService.GetRepositoryByName(username, repoName)
	if errGetRep != nil {
		utils.RespondErrorString(c, "Could not find repository: "+repoName, http.StatusNotFound)
		return
	}

	// Get image from repository
	imageTag, errGetTag := i.ImageService.GetImageTagByRepoAndTag(repo.Id, imageName)
	if errGetTag != nil {
		utils.RespondErrorString(c, "Could not find: "+username+"/"+repoName+":"+imageName, http.StatusNotFound)
		return
	}

	// Get the download_url
	URL, errGetURL := i.StorageService.GetDownloadURL(imageTag.FileKey)
	if errGetURL != nil {
		utils.RespondError(c, errGetURL, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"download_url": URL, "data": imageTag})

}

func (i *ImageController) GetImageTagsForRepoName(c *gin.Context) {
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")

	repo, errGetRepo := i.RepositoryService.GetRepositoryByName(username, repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Repository not found", 404)
		return
	}
	imageTags, err := i.ImageService.GetImageTagsForRepo(repo.Id)
	if err != nil {
		utils.RespondSQLError(c, err)
		return
	}

	c.JSON(200, gin.H{"images": &imageTags})
}
