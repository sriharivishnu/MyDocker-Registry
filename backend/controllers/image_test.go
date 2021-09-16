package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mocks "github.com/sriharivishnu/shopify-challenge/mocks/layers"
	serviceMocks "github.com/sriharivishnu/shopify-challenge/mocks/services"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestImagePull(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	t.Run("Push Image Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{
			{
				Key:   "user_id",
				Value: "123",
			},
			{
				Key:   "repo_id",
				Value: "456",
			},
			{
				Key:   "image_id",
				Value: "srihari",
			},
		}

		dummyRepo := models.Repository{Id: "456"}
		mockRepoService := mocks.RepositoryLayer{}
		mockRepoService.On("GetRepositoryByName", "123", "456").Return(dummyRepo, nil)

		dummyImageTag := models.ImageTag{
			Id:           "121314",
			RepositoryId: "456",
			Tag:          "srihari",
			FileKey:      "test/file/key.tar.gz",
		}
		mockImageService := mocks.ImageLayer{}
		mockImageService.On("GetImageTagByRepoAndTag", "456", "srihari").Return(dummyImageTag, nil)

		mockStorageService := serviceMocks.Storage{}
		mockStorageService.On("GetDownloadURL", "test/file/key.tar.gz").Return("http://url.com", nil)

		imageController := ImageController{
			RepositoryService: &mockRepoService,
			ImageService:      &mockImageService,
			StorageService:    &mockStorageService,
		}

		imageController.PullImage(ctx)

		expected, _ := json.Marshal(gin.H{
			"download_url": "http://url.com",
			"data":         &dummyImageTag,
		})

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expected, w.Body.Bytes())
	})
}
