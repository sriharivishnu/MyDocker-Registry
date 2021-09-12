package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/services"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type ImageController struct{}

func (*ImageController) GetImageTagsForRepoName(repoName string) ([]models.ImageTag, error) {
	return nil, nil
}

type UploadURLPayload struct {
	RepoName string `json:"repository_name"`
	ImageTag string `json:"image_tag"`
}

func (*ImageController) GetUploadURL(c *gin.Context) {
	payload := UploadURLPayload{}
	errInputFormat := c.BindJSON(&payload)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	repo := models.Repository{}
	errGetRep := repo.GetRepositoryByName(payload.RepoName)
	if errGetRep != nil {
		utils.RespondErrorString(c, "Could not find repository: "+payload.RepoName, http.StatusNotFound)
		return
	}

	storage := services.S3{}
	URL, err := storage.GetUploadURL(user.Username, payload.RepoName, payload.ImageTag)

	if err != nil {
		utils.RespondError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"url": URL})
}

func (*ImageController) CreateImageTag(c *gin.Context) {}
