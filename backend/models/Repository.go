package models

import (
	"fmt"
	"strings"
	"time"
)

type Repository struct {
	Id          string    `json:"id,omitempty" db:"id"`
	OwnerId     string    `json:"owner_id,omitempty" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

func (repo Repository) Validate() error {
	if len(repo.Name) == 0 {
		return fmt.Errorf("Repository name must not be empty")
	}
	if strings.Contains(repo.Name, "/") || strings.Contains(repo.Name, "\\") {
		return fmt.Errorf("Repository name contains unknown characters: %s", repo.Name)
	}

	return nil
}
