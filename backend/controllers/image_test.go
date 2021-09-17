package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	utils "github.com/sriharivishnu/shopify-challenge/mocks/helpers"
	mocks "github.com/sriharivishnu/shopify-challenge/mocks/layers"
	serviceMocks "github.com/sriharivishnu/shopify-challenge/mocks/services"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestImagePush(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	params := []gin.Param{
		{
			Key:   "user_id",
			Value: "username",
		},
		{
			Key:   "repo_id",
			Value: "456",
		},
	}
	dummyRepo := models.Repository{Id: "456", OwnerId: "user_id", Name: "repoName"}
	dummyImageTag := models.ImageTag{
		Id:           "121314",
		RepositoryId: "456",
		Tag:          "srihari",
		FileKey:      "test/file/key.tar.gz",
	}

	dummyUser := models.User{Id: "user_id", Username: "username"}

	t.Run("Push Image Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = params
		ctx.Set("user", dummyUser)

		utils.MockJsonPost(ctx, gin.H{"tag": "srihari", "description": "sample description"})

		mockRepoService := mocks.RepositoryLayer{}
		mockRepoService.On("GetRepositoryByName", dummyUser.Username, "456").Return(dummyRepo, nil)

		mockImageService := mocks.ImageLayer{}
		mockImageService.On("Create", "456", "srihari", "sample description", "username/repoName/srihari.tar.gz").Return(dummyImageTag, nil)

		mockStorageService := serviceMocks.Storage{}
		mockStorageService.On("GetUploadURL", dummyUser.Username, "456", "srihari").Return("http://url.com", nil)

		imageController := ImageController{
			RepositoryService: &mockRepoService,
			ImageService:      &mockImageService,
			StorageService:    &mockStorageService,
		}

		imageController.PushImage(ctx)

		expected, _ := json.Marshal(gin.H{
			"message":    "Created image successfully",
			"upload_url": "http://url.com",
			"id":         dummyImageTag.Id,
		})

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expected, w.Body.Bytes())
	})
}

func TestImagePull(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)

	params := []gin.Param{
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
	dummyImageTag := models.ImageTag{
		Id:           "121314",
		RepositoryId: "456",
		Tag:          "srihari",
		FileKey:      "test/file/key.tar.gz",
	}
	t.Run("Pull Image Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = params

		mockRepoService := mocks.RepositoryLayer{}
		mockRepoService.On("GetRepositoryByName", "123", "456").Return(dummyRepo, nil)

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

	t.Run("Pull Image Error", func(t *testing.T) {
		testcases := []struct {
			getRepoError  error
			getImageError error
			getURLError   error
			expected      string
			code          int
		}{
			{
				sql.ErrNoRows, nil, nil, "Could not find repository: 456", 404,
			},
			{
				nil, sql.ErrNoRows, nil, "Could not find: 123/456:srihari", 404,
			},
			{
				nil, nil, errors.New("error fetching URL"), "error fetching URL", 500,
			},
		}

		for _, test := range testcases {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = params

			mockRepoService := mocks.RepositoryLayer{}
			mockRepoService.On("GetRepositoryByName", "123", "456").Return(models.Repository{}, test.getRepoError)

			mockImageService := mocks.ImageLayer{}
			mockImageService.On("GetImageTagByRepoAndTag", "456", "srihari").Return(models.ImageTag{}, test.getImageError)

			mockStorageService := serviceMocks.Storage{}
			mockStorageService.On("GetDownloadURL", "test/file/key.tar.gz").Return("http://url.com", test.getURLError)

			imageController := ImageController{
				RepositoryService: &mockRepoService,
				ImageService:      &mockImageService,
				StorageService:    &mockStorageService,
			}

			imageController.PullImage(ctx)

			expected, _ := json.Marshal(gin.H{
				"error": test.expected,
			})

			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, expected, w.Body.Bytes())
		}

	})
}
