package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (controller *AuthController) SignUp(c *gin.Context) {
	var user models.User
	errInputFormat := c.BindJSON(&user)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errInputFormat != nil {
		utils.RespondError(c, errHash, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	errUserCreate := user.Create()
	if errUserCreate != nil {
		utils.RespondError(c, errUserCreate, http.StatusInternalServerError)
		return
	}

	token, errToken := createToken(user)
	if errToken != nil {
		utils.RespondError(c, errToken, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Signed up successfully", "token": token})

}

func (controller *AuthController) SignIn(c *gin.Context) {
	var userPayload models.User
	errInputFormat := c.BindJSON(&userPayload)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	var userDB models.User
	errGetUser := userDB.GetByUsername(userPayload.Username)
	if errGetUser != nil {
		if errGetUser == sql.ErrNoRows {
			utils.RespondErrorString(c, "Unauthorized", http.StatusUnauthorized)
			return
		}
		utils.RespondError(c, errGetUser, http.StatusInternalServerError)
		return
	}

	if errHash := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(userPayload.Password)); errHash != nil {
		utils.RespondErrorString(c, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, errToken := createToken(userDB)
	if errToken != nil {
		utils.RespondError(c, errToken, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Signed in successfully", "token": token})

}

func createToken(user models.User) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := at.SignedString([]byte(config.Config.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return token, nil

}
