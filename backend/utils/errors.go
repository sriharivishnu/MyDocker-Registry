package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func RespondSQLError(c *gin.Context, err error) {
	log.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		RespondError(c, err, 500)
		return
	}
	if sqlErr.Number == 1062 || sqlErr.Number == 1169 {
		RespondErrorString(c, "This Resource Already Exists!", http.StatusConflict)
		return
	}
	RespondErrorString(c, "Internal Server Error", 500)
}

func RespondError(c *gin.Context, err error, errorCode int) {
	c.JSON(errorCode, gin.H{"error": err.Error()})
	log.Printf("\x1b[31;1m%s\x1b[0m\n", err)
}

func RespondErrorString(c *gin.Context, message string, errorCode int) {
	c.JSON(errorCode, gin.H{"error": message})
	log.Printf("\x1b[31;1m%s\x1b[0m\n", message)
}
