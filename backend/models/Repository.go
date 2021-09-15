package models

import (
	"log"
	"time"

	"github.com/pkg/errors"
	db "github.com/sriharivishnu/shopify-challenge/services"
)

type Repository struct {
	Id          string    `json:"id,omitempty" db:"id"`
	OwnerId     string    `json:"owner_id,omitempty" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IRepository interface {
	Create() error
	GetRepositoryByName(name string) error
	GetRepositoriesForUser(ownerId string) ([]Repository, error)
	Search(query string, limit int, offset int) ([]Repository, error)
}

func (repo *Repository) Create() error {
	tx := db.DbConn.MustBegin()
	err := tx.Get(repo, "INSERT INTO repository (name, description, owner_id) VALUES (?, ?, ?) returning *;", repo.Name, repo.Description, repo.OwnerId)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "create repository error")
	}
	tx.Commit()
	return nil
}

func (repo *Repository) GetRepositoryByName(username string, reponame string) error {
	log.Println(username, reponame)
	sql := `select r.* from repository r
				INNER JOIN user u on r.owner_id = u.id
				where u.username = ? and r.name = ?;`
	err := db.DbConn.Get(repo, sql, username, reponame)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetRepositoriesForUser(ownerId string) ([]Repository, error) {
	sql := `select r.* from repository r
		inner join user u on r.owner_id = u.id
		where r.owner_id = ? or u.username = ?`
	repos := []Repository{}
	err := db.DbConn.Select(&repos, sql, ownerId, ownerId)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (repo *Repository) Search(query string, limit int, offset int) ([]Repository, error) {
	return nil, nil
}
