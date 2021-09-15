package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type RepositoryController struct{}

func (*RepositoryController) Create(c *gin.Context) {
	var repo models.Repository
	errInputFormat := c.BindJSON(&repo)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	repo.OwnerId = user.Id
	errCreate := repo.Create()
	if errCreate != nil {
		utils.RespondSQLError(c, errCreate)
		return
	}

	c.JSON(200, &repo)

}

func (*RepositoryController) GetForUser(c *gin.Context) {
	user_id, _ := c.Params.Get("user_id")

	m := models.Repository{}
	repos, err := m.GetRepositoriesForUser(user_id)
	if err != nil {
		utils.RespondSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"repositories": &repos})
}

func (*RepositoryController) Search(c *gin.Context) {
}
