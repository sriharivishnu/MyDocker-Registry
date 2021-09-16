package models

import (
	"time"
)

type ImageTag struct {
	Id           string    `json:"id" db:"id"`
	RepositoryId string    `json:"repository_id" db:"repository_id"`
	Description  string    `json:"description,omitempty" db:"description"`
	Tag          string    `json:"tag,omitempty" db:"tag"`
	FileKey      string    `json:"file_key,omitempty" db:"file_key"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
}
