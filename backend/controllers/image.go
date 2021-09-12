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

func (*ImageController) GetUploadURL(c *gin.Context) {
	repoName, foundRepo := c.Params.Get("repo")
	imageTag, foundTag := c.Params.Get("tag")

	if !foundRepo || !foundTag {
		utils.RespondErrorString(c, "Invalid parameters", http.StatusBadRequest)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	repo := models.Repository{}
	errGetRepo := repo.GetRepositoryByName(repoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Could not find repository: "+repoName, http.StatusNotFound)
		return
	}
	if user.Id != repo.OwnerId {
		utils.RespondErrorString(c, "User is not allowed to upload for this repository", http.StatusForbidden)
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

type createImageTagPayload struct {
	RepoName    string `json:"repository_name"`
	ImageTag    string `json:"image_tag"`
	Description string `json:"description,omitempty"`
}

func (*ImageController) CreateImageTag(c *gin.Context) {
	var payload createImageTagPayload
	errInputFormat := c.BindJSON(&payload)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	repo := models.Repository{}
	errGetRepo := repo.GetRepositoryByName(payload.RepoName)
	if errGetRepo != nil {
		utils.RespondErrorString(c, "Could not find repository: "+payload.RepoName, http.StatusNotFound)
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
	errCreate := imageTag.Create()
	if errCreate != nil {
		utils.RespondSQLError(c, errCreate)
		return
	}

	c.JSON(200, gin.H{"message": "Created image successfully", "id": imageTag.Id})

}

func (*ImageController) GetImage(c *gin.Context) {
	repoName, foundRepo := c.Params.Get("repo")
	tagName, foundTag := c.Params.Get("tag")
	if !foundRepo || !foundTag {
		utils.RespondErrorString(c, "Invalid parameters", http.StatusBadRequest)
		return
	}

	repo := models.Repository{}
	errGetRep := repo.GetRepositoryByName(repoName)
	if errGetRep != nil {
		utils.RespondErrorString(c, "Could not find repository: "+repoName, http.StatusNotFound)
		return
	}

	imageTag := models.ImageTag{}
	errGetTag := imageTag.GetImageTagByRepoAndTag(repo.Id, tagName)
	if errGetTag != nil {
		utils.RespondErrorString(c, "Could not find: "+repoName+":"+tagName, http.StatusNotFound)
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
