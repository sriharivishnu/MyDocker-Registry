package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func CheckSQLError(err error) bool {
	err, ok := err.(*mysql.MySQLError)
	return ok
}

func RespondSQLError(c *gin.Context, err error) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return
	}
	if sqlErr.Number == 1062 || sqlErr.Number == 1169 {
		RespondErrorString(c, "This Resource Already Exists!", http.StatusConflict)
		return
	}
	RespondError(c, err, 500)
}

func RespondError(c *gin.Context, err error, errorCode int) {
	c.JSON(errorCode, gin.H{"error": err.Error()})
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
}

func RespondErrorString(c *gin.Context, message string, errorCode int) {
	c.JSON(errorCode, gin.H{"error": message})
}
