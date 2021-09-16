package layers

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/models"
	db "github.com/sriharivishnu/shopify-challenge/services"
)

type UserLayer interface {
	Create(username, password string) (models.User, error)
	GetByUsername(username string) (models.User, error)
	GetByID(userId string) (models.User, error)
	CreateToken(user models.User) (string, error)
}

type UserService struct{}

func (service *UserService) Create(username, password string) (models.User, error) {
	user := models.User{}
	tx := db.DbConn.MustBegin()
	err := tx.Get(&user, "INSERT INTO user (username, password) VALUES (?, ?) returning *;", username, password)
	if err != nil {
		tx.Rollback()
		return user, err
	}
	tx.Commit()
	return user, nil
}

func (service *UserService) GetByUsername(username string) (models.User, error) {
	user := models.User{}
	err := db.DbConn.Get(&user, "select * from user where username = ?", username)
	return user, err
}

func (service *UserService) GetByID(userId string) (models.User, error) {
	user := models.User{}
	err := db.DbConn.Get(&user, "select * from user where id = ?", userId)
	return user, err
}

func (service *UserService) CreateToken(user models.User) (string, error) {
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
