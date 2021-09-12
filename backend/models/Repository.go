package models

import "time"

type Repository struct {
	Id          string    `json:"id" db:"id"`
	OwnerId     string    `json:"owner_id" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IRepository interface {
	GetRepositoryById(id string) error
	GetRepositoriesForUser(ownerId string) ([]Repository, error)
	Search(query string, limit int, offset int) ([]Repository, error)
}
