package services

import (
	"github.com/pkg/errors"
	db "github.com/sriharivishnu/shopify-challenge/external"
	"github.com/sriharivishnu/shopify-challenge/models"
)

type RepositoryLayer interface {
	Create(name, description, ownerId string) (models.Repository, error)
	GetRepositoryByName(username string, reponame string) (models.Repository, error)
	GetRepositoriesForUser(ownerId string) ([]models.Repository, error)
	Search(query string, limit int, offset int) ([]models.SearchResult, error)
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

func (service *RepositoryService) Search(query string, limit int, offset int) ([]models.SearchResult, error) {
	repos := []models.SearchResult{}
	sql := `
		select 	r.id as id,
				r.name as name,
				r.description as description,
				u.username as username,
				r.created_at as created_at,
				COALESCE(counts.num_tags, 0) as num_tags
		from repository r
				inner join user u on r.owner_id = u.id
				left join (
					select count(*) as num_tags, repository_id from image_tag group by repository_id
				) as counts on counts.repository_id = r.id
				where MATCH(r.name) AGAINST(?) limit ? offset ?;`
	query = "*" + query + "*"
	err := db.DbConn.Select(&repos, sql, query, limit, offset)
	return repos, err
}
