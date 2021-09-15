package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/services"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type ImageController struct{}

func (*ImageController) GetUploadURL(c *gin.Context) {
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")
	imageTag, _ := c.Params.Get("image_id")

	if len(imageTag) == 0 {
		utils.RespondErrorString(c, "Please tag your image before pushing!", http.StatusNotFound)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	repo := models.Repository{}
	errGetRepo := repo.GetRepositoryByName(username, repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Repository not found", http.StatusNotFound)
		return
	}
	if repo.OwnerId != user.Id {
		utils.RespondErrorString(c, "User is not authorized to push to repository", http.StatusForbidden)
		return
	}

	storage := services.S3{}
	URL, err := storage.GetUploadURL(user.Username, repoName, imageTag)

	if err != nil {
		utils.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"upload_url": URL})
}

func (*ImageController) CreateImageTag(c *gin.Context) {
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

	repo := models.Repository{}
	errGetRepo := repo.GetRepositoryByName(username, repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Could not find repository: "+username+"/"+repoName, http.StatusNotFound)
		return
	}
	if user.Id != repo.OwnerId {
		utils.RespondErrorString(c, "User is not allowed to upload for this repository", http.StatusForbidden)
		return
	}

	key := utils.CreateFileKey(user.Username, repo.Name, payload.ImageTag)
	imageTag := models.ImageTag{
		RepositoryId: repo.Id,
		Description:  payload.Description,
		Tag:          payload.ImageTag,
		FileKey:      key,
	}
	log.Println("tag" + payload.ImageTag)
	errCreate := imageTag.Create()
	if errCreate != nil {
		utils.RespondSQLError(c, errCreate)
		return
	}

	c.JSON(200, gin.H{"message": "Created image successfully", "id": imageTag.Id})

}

func (*ImageController) GetImage(c *gin.Context) {
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")
	imageName, _ := c.Params.Get("image_id")

	repo := models.Repository{}
	errGetRep := repo.GetRepositoryByName(username, repoName)
	if errGetRep != nil {
		utils.RespondErrorString(c, "Could not find repository: "+repoName, http.StatusNotFound)
		return
	}

	imageTag := models.ImageTag{}
	var errGetTag error
	if imageName == "" || imageName == "latest" {
		errGetTag = imageTag.GetLatestImageTag(repo.Id)
	} else {
		errGetTag = imageTag.GetImageTagByRepoAndTag(repo.Id, imageName)
	}
	if errGetTag != nil {
		utils.RespondErrorString(c, "Could not find: "+username+"/"+repoName+":"+imageName, http.StatusNotFound)
		return
	}

	storage := services.S3{}
	URL, err := storage.GetDownloadURL(imageTag.FileKey)

	if err != nil {
		utils.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"download_url": URL, "data": imageTag})

}

func (*ImageController) GetImageTagsForRepoName(c *gin.Context) {
	username, _ := c.Params.Get("user_id")
	repoName, _ := c.Params.Get("repo_id")
	repo := models.Repository{}

	errGetRepo := repo.GetRepositoryByName(username, repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Repository not found", 404)
		return
	}
	getTags := models.ImageTag{}
	imageTags, err := getTags.GetImageTagsForRepo(repo.Id)
	if err != nil {
		utils.RespondSQLError(c, err)
		return
	}

	c.JSON(200, gin.H{"images": &imageTags})
}
