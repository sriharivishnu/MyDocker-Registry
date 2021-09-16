package layers

import (
	"github.com/pkg/errors"
	"github.com/sriharivishnu/shopify-challenge/models"
	db "github.com/sriharivishnu/shopify-challenge/services"
)

type RepositoryLayer interface {
	Create(name, description, ownerId string) (models.Repository, error)
	GetRepositoryByName(username string, reponame string) (models.Repository, error)
	GetRepositoriesForUser(ownerId string) ([]models.Repository, error)
	Search(query string, limit int, offset int) ([]models.Repository, error)
}

type RepositoryService struct{}

func (service *RepositoryService) Create(name, description, ownerId string) (models.Repository, error) {
	repo := models.Repository{}
	tx := db.DbConn.MustBegin()
	err := tx.Get(&repo, "INSERT INTO repository (name, description, owner_id) VALUES (?, ?, ?) returning *;", name, description, ownerId)
	if err != nil {
		tx.Rollback()
		return repo, errors.Wrap(err, "create repository error")
	}
	tx.Commit()
	return repo, nil
}

func (service *RepositoryService) GetRepositoryByName(username string, reponame string) (models.Repository, error) {
	repo := models.Repository{}
	sql := `select r.* from repository r
				INNER JOIN user u on r.owner_id = u.id
				where u.username = ? and r.name = ?;`
	err := db.DbConn.Get(&repo, sql, username, reponame)
	return repo, err
}

func (service *RepositoryService) GetRepositoriesForUser(ownerId string) ([]models.Repository, error) {
	repos := []models.Repository{}
	sql := `select r.* from repository r
		inner join user u on r.owner_id = u.id
		where r.owner_id = ? or u.username = ?`
	err := db.DbConn.Select(&repos, sql, ownerId, ownerId)
	return repos, err
}

func (service *RepositoryService) Search(query string, limit int, offset int) ([]models.Repository, error) {
	return nil, nil
}
