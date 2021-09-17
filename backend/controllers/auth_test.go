package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	utils "github.com/sriharivishnu/shopify-challenge/mocks/helpers"
	mocks "github.com/sriharivishnu/shopify-challenge/mocks/layers"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpSuccess(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	t.Run("Sign Up Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		utils.MockJsonPost(ctx, gin.H{"username": "srihari", "password": "testpassword"})

		dummyUser := models.User{Username: "srihari", Password: "testpassword", Id: "123", CreatedAt: time.Now()}

		mockUserService := mocks.UserLayer{}
		mockUserService.On("Create", "srihari", mock.AnythingOfType("string")).Return(dummyUser, nil)
		mockUserService.On("CreateToken", dummyUser).Return("dummy_token", nil)

		authController := AuthController{
			UserService: &mockUserService,
		}
		authController.SignUp(ctx)

		expected, _ := json.Marshal(gin.H{
			"message": "Signed up successfully",
			"token":   "dummy_token",
		})

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expected, w.Body.Bytes())
	})
}

func TestSignUpError(t *testing.T) {
	testcases := []struct {
		testName    string
		username    string
		password    string
		response    string
		createError error
		code        int
	}{
		{
			"InvalidUsername", "sri", "testpassword", "username must be at least 5 characters in length", nil, 400,
		},
		{
			"InvalidPassword", "srihari", "test", "password must be at least 6 characters in length", nil, 400,
		},
		{
			"UserExists", "srihari", "test123", "This Resource Already Exists!", &mysql.MySQLError{Number: 1062}, 409,
		},
		{
			"CreateError", "srihari", "test123", "sql error", errors.New("sql error"), 500,
		},
	}
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	for _, test := range testcases {

		t.Run(test.testName, func(t *testing.T) {

			for _, test := range testcases {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				utils.MockJsonPost(ctx, gin.H{"username": test.username, "password": test.password})

				mockUserService := mocks.UserLayer{}
				mockUserService.On("Create", test.username, mock.AnythingOfType("string")).Return(models.User{}, test.createError)

				authController := AuthController{
					UserService: &mockUserService,
				}
				authController.SignUp(ctx)

				assert.Equal(t, test.code, w.Code)
				expected, _ := json.Marshal(gin.H{
					"error": test.response,
				})
				assert.Equal(t, expected, w.Body.Bytes())
			}
		})
	}
}

func TestSignInSuccess(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	t.Run("Sign In Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		utils.MockJsonPost(ctx, gin.H{"username": "srihari", "password": "testpassword"})

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), 8)

		dummyUser := models.User{Username: "srihari", Password: string(hashedPassword), Id: "123", CreatedAt: time.Now()}

		mockUserService := mocks.UserLayer{}
		mockUserService.On("GetByUsername", "srihari", mock.AnythingOfType("string")).Return(dummyUser, nil)
		mockUserService.On("CreateToken", dummyUser).Return("dummy_token", nil)

		authController := AuthController{
			UserService: &mockUserService,
		}
		authController.SignIn(ctx)

		assert.Equal(t, 200, w.Code)
		expected, _ := json.Marshal(gin.H{
			"message": "Signed in successfully", "token": "dummy_token",
		})
		assert.Equal(t, expected, w.Body.Bytes())
	})

}

func TestSignInError(t *testing.T) {
	testcases := []struct {
		testName     string
		username     string
		password     string
		dbPassword   string
		getUserError error
		response     string
		code         int
	}{
		{
			"testSQLError", "srihari", "testpassword", "testpassword", errors.New("sql error"), "sql error", 500,
		},
		{
			"InvalidSignIn", "srihari", "testpassword", "notright", nil, "Username or password is incorrect. Please check your login details and try again.", 401,
		},
		{
			"UserNotFound", "srihari", "testpassword", "testpassword", sql.ErrNoRows, "Username or password is incorrect. Please check your login details and try again.", 401,
		},
	}
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.TestMode)
	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			utils.MockJsonPost(ctx, gin.H{"username": test.username, "password": test.password})

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(test.dbPassword), 8)

			dummyUser := models.User{Username: test.username, Password: string(hashedPassword), Id: "123", CreatedAt: time.Now()}

			mockUserService := mocks.UserLayer{}
			mockUserService.On("GetByUsername", test.username).Return(dummyUser, test.getUserError)

			authController := AuthController{
				UserService: &mockUserService,
			}
			authController.SignIn(ctx)

			assert.Equal(t, test.code, w.Code)
			expected, _ := json.Marshal(gin.H{
				"error": test.response,
			})
			assert.Equal(t, expected, w.Body.Bytes())
		})
	}
}
