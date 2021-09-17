package models

import (
	"fmt"
	"time"

	"github.com/sriharivishnu/shopify-challenge/utils"
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
	if !utils.IsValidName(repo.Name) {
		return fmt.Errorf("Repository name contains invalid characters. Please only use letters, numbers, and/or -,_")
	}
	return nil
}
