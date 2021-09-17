package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	utils "github.com/sriharivishnu/shopify-challenge/mocks/helpers"
	mocks "github.com/sriharivishnu/shopify-challenge/mocks/services"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreateSuccess(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	params := []gin.Param{
		{
			Key:   "user_id",
			Value: "username",
		},
	}
	dummyRepo := models.Repository{Id: "456", OwnerId: "user_id", Name: "repoName", Description: "sample description"}

	dummyUser := models.User{Id: "user_id", Username: "username"}

	t.Run("RepositoryCreateSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = params
		ctx.Set("user", dummyUser)

		utils.MockJsonPost(ctx, gin.H{"name": "repoName", "description": "sample description"})

		mockRepoService := mocks.RepositoryLayer{}
		mockRepoService.On("Create", dummyRepo.Name, dummyRepo.Description, dummyRepo.OwnerId).Return(dummyRepo, nil)

		repoController := RepositoryController{
			RepositoryService: &mockRepoService,
		}

		repoController.Create(ctx)

		assert.Equal(t, 200, w.Code)
		expected, _ := json.Marshal(gin.H{"message": "Created repository username/repoName successfully", "id": "456"})
		assert.Equal(t, expected, w.Body.Bytes())
	})

}

func TestRepositoryCreateError(t *testing.T) {
	testcases := []struct {
		testName    string
		repoName    string
		description string
		response    string
		createError error
		code        int
	}{
		{
			"EmptyRepoName", "", "test", "Repository name must not be empty", nil, 400,
		},
		{
			"InvalidRepoName", "/\\", "test", "Repository name contains invalid characters. Please only use letters, numbers, and/or -,_", nil, 400,
		},
		{
			"RepoCreateError", "test123", "desc", "Test SQL error", errors.New("Test SQL error"), 500,
		},
	}
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	params := []gin.Param{
		{
			Key:   "user_id",
			Value: "username",
		},
	}

	dummyUser := models.User{Id: "user_id", Username: "username"}

	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = params
			ctx.Set("user", dummyUser)

			utils.MockJsonPost(ctx, gin.H{"name": test.repoName, "description": test.description})

			mockRepoService := mocks.RepositoryLayer{}
			mockRepoService.On("Create", test.repoName, test.description, dummyUser.Id).Return(models.Repository{Id: "test"}, test.createError)

			repoController := RepositoryController{
				RepositoryService: &mockRepoService,
			}

			repoController.Create(ctx)

			assert.Equal(t, test.code, w.Code)
			expected, _ := json.Marshal(gin.H{"error": test.response})
			assert.Equal(t, expected, w.Body.Bytes())
		})
	}

}

func TestRepositoryGetForUser(t *testing.T) {
	dummyRepositories := []models.Repository{
		{
			Id:          "123",
			OwnerId:     "456",
			Name:        "foo",
			Description: "test description",
		},
		{
			Id:          "789",
			OwnerId:     "456",
			Name:        "bar",
			Description: "test description 2",
		},
		{
			Id:          "101",
			OwnerId:     "456",
			Name:        "cool",
			Description: "",
		},
	}

	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	params := []gin.Param{
		{
			Key:   "user_id",
			Value: "username",
		},
	}

	t.Run("GetRepositoriesSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = params

		mockRepoService := mocks.RepositoryLayer{}
		mockRepoService.On("GetRepositoriesForUser", "username").Return(dummyRepositories, nil)

		repoController := RepositoryController{
			RepositoryService: &mockRepoService,
		}

		repoController.GetForUser(ctx)

		assert.Equal(t, 200, w.Code)
		expected, _ := json.Marshal(gin.H{"repositories": dummyRepositories})
		assert.Equal(t, expected, w.Body.Bytes())
	})

}
