package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/layers"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type RepositoryController struct {
	RepositoryService layers.RepositoryLayer
}

func (r *RepositoryController) Create(c *gin.Context) {
	curUser, _ := c.Get("user")
	user := curUser.(models.User)

	// Input validation
	var repo models.Repository
	if errInputFormat := c.BindJSON(&repo); errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}
	if errValidation := repo.Validate(); errValidation != nil {
		utils.RespondError(c, errValidation, http.StatusBadRequest)
		return
	}

	// Create repository in DB
	repo.OwnerId = user.Id
	repo, errCreate := r.RepositoryService.Create(repo.Name, repo.Description, repo.OwnerId)
	if errCreate != nil {
		utils.RespondSQLError(c, errCreate)
		return
	}

	c.JSON(200, &repo)

}

func (r *RepositoryController) GetForUser(c *gin.Context) {
	user_id, _ := c.Params.Get("user_id")

	repos, err := r.RepositoryService.GetRepositoriesForUser(user_id)
	if err != nil {
		utils.RespondSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"repositories": &repos})
}

func (r *RepositoryController) Search(c *gin.Context) {
	params := c.Request.URL.Query()
	query := params.Get("query")
	if query == "" {
		utils.RespondErrorString(c, "Query cannot be empty", 400)
		return
	}

	offsetStr := params.Get("offset")
	offset, offsetErr := strconv.Atoi(offsetStr)
	if offsetErr != nil {
		offset = 0
	}
	repos, err := r.RepositoryService.Search(query, 10, offset)
	if err != nil {
		utils.RespondSQLError(c, err)
		return
	}
	c.JSON(200, gin.H{"results": &repos})
}
