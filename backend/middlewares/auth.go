package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/models"
)

func getUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.JWT_SECRET), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || err != nil {
		return nil, err
	}

	user_id := claims["user_id"].(string)
	user := &models.User{}
	errGetUser := user.GetById(user_id)
	if errGetUser != nil {
		return nil, errGetUser
	}
	return user, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		token := strings.TrimSpace(auth[7:])
		user, err := getUserFromToken(token)
		if err != nil {
			c.AbortWithError(403, err)
			return
		}
		if user == nil {
			c.AbortWithError(500, fmt.Errorf("cannot find user"))
			return
		}
		c.Set("user", *user)
		c.Next()
	}
}
