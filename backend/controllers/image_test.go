package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	externalMocks "github.com/sriharivishnu/shopify-challenge/mocks/external"
	utils "github.com/sriharivishnu/shopify-challenge/mocks/helpers"
	mocks "github.com/sriharivishnu/shopify-challenge/mocks/services"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestImagePushSuccess(t *testing.T) {
	testcases := []struct {
		testName    string
		tag         string
		expectedTag string
	}{
		{
			"TestSuccess", "srihari", "srihari",
		},
		{
			"TestLatestSuccess", "", "latest",
		},
	}
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

	dummyUser := models.User{Id: "user_id", Username: "username"}

	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {

			dummyImageTag := models.ImageTag{
				Id:           "121314",
				RepositoryId: "456",
				Tag:          test.tag,
				FileKey:      "test/file/key.tar.gz",
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = params
			ctx.Set("user", dummyUser)

			utils.MockJsonPost(ctx, gin.H{"tag": test.tag, "description": "sample description"})

			mockRepoService := mocks.RepositoryLayer{}
			mockRepoService.On("GetRepositoryByName", dummyUser.Username, "456").Return(dummyRepo, nil)

			mockImageService := mocks.ImageLayer{}
			fileKey := fmt.Sprintf("username/repoName/%s.tar.gz", test.expectedTag)
			mockImageService.On("Create", "456", test.expectedTag, "sample description", fileKey).Return(dummyImageTag, nil)

			mockStorageService := externalMocks.Storage{}
			mockStorageService.On("GetUploadURL", dummyUser.Username, "456", test.expectedTag).Return("http://url.com", nil)

			imageController := ImageController{
				RepositoryService: &mockRepoService,
				ImageService:      &mockImageService,
				StorageService:    &mockStorageService,
			}

			imageController.PushImage(ctx)

			expected, _ := json.Marshal(gin.H{
				"message":    fmt.Sprintf("Created image username/repoName:%s successfully", test.expectedTag),
				"upload_url": "http://url.com",
				"id":         dummyImageTag.Id,
			})

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, expected, w.Body.Bytes())

		})
	}

}

func TestImagePushError(t *testing.T) {
	testcases := []struct {
		testName       string
		getRepoError   error
		createError    error
		getUploadError error
		code           int
		response       string
	}{
		{
			"RepoNotFound", sql.ErrNoRows, nil, nil, 404, "Could not find repository: username/456",
		},
		{
			"ImageExists", nil, &mysql.MySQLError{Number: 1062}, nil, 409, "This Resource Already Exists!",
		},
		{
			"UploadError", nil, nil, errors.New("Get upload error"), 500, "Get upload error",
		},
	}
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

	dummyUser := models.User{Id: "user_id", Username: "username"}
	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = params
			ctx.Set("user", dummyUser)

			utils.MockJsonPost(ctx, gin.H{"tag": "srihari", "description": "sample description"})

			mockRepoService := mocks.RepositoryLayer{}
			mockRepoService.On("GetRepositoryByName", dummyUser.Username, "456").Return(dummyRepo, test.getRepoError)

			mockImageService := mocks.ImageLayer{}
			mockImageService.On("Create", "456", "srihari", "sample description", "username/repoName/srihari.tar.gz").Return(models.ImageTag{Id: "test"}, test.createError)

			mockStorageService := externalMocks.Storage{}
			mockStorageService.On("GetUploadURL", dummyUser.Username, "456", "srihari").Return("http://url.com", test.getUploadError)

			imageController := ImageController{
				RepositoryService: &mockRepoService,
				ImageService:      &mockImageService,
				StorageService:    &mockStorageService,
			}

			imageController.PushImage(ctx)

			expected, _ := json.Marshal(gin.H{
				"error": test.response,
			})

			assert.Equal(t, test.code, w.Code)
			assert.Equal(t, expected, w.Body.Bytes())
		})
	}

}

func TestImagePullSuccess(t *testing.T) {
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
	dummyRepo := models.Repository{Id: "456", Name: "repo_name"}
	dummyImageTag := models.ImageTag{
		Id:           "121314",
		RepositoryId: dummyRepo.Id,
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

		mockStorageService := externalMocks.Storage{}
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

func TestImagePullError(t *testing.T) {
	testcases := []struct {
		testName      string
		getRepoError  error
		getImageError error
		getURLError   error
		expected      string
		code          int
	}{
		{
			"NoRepository", sql.ErrNoRows, nil, nil, "Could not find repository: 456", 404,
		},
		{
			"ImageNotFound", nil, sql.ErrNoRows, nil, "Could not find: 123/456:srihari", 404,
		},
		{
			"URLError", nil, nil, errors.New("error fetching URL"), "error fetching URL", 500,
		},
	}
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
	dummyRepo := models.Repository{Id: "456", Name: "repo_name"}
	dummyImageTag := models.ImageTag{
		Id:           "121314",
		RepositoryId: dummyRepo.Id,
		Tag:          "srihari",
		FileKey:      "test/file/key.tar.gz",
	}
	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = params

			mockRepoService := mocks.RepositoryLayer{}
			mockRepoService.On("GetRepositoryByName", "123", "456").Return(dummyRepo, test.getRepoError)

			mockImageService := mocks.ImageLayer{}
			mockImageService.On("GetImageTagByRepoAndTag", "456", "srihari").Return(dummyImageTag, test.getImageError)

			mockStorageService := externalMocks.Storage{}
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

		})
	}
}
