package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/layers"
	"github.com/sriharivishnu/shopify-challenge/models"
)

func getUserFromToken(tokenString string, userService layers.UserLayer) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.JWT_SECRET), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || err != nil {
		return nil, fmt.Errorf("invalid token claims: token may be expired, please try logging in again")
	}

	user_id := claims["user_id"].(string)
	user, errGetUser := userService.GetByID(user_id)
	if errGetUser != nil {
		return nil, errGetUser
	}
	return &user, nil
}

func AuthMiddleware(userService layers.UserLayer) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if len(auth) <= 7 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid Authorization header"})
			return
		}
		token := strings.TrimSpace(auth[7:])
		user, err := getUserFromToken(token, userService)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "User not found"})
			return
		}
		c.Set("user", *user)
		c.Next()
	}
}
