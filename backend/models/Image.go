package models

import "time"

type ImageTag struct {
	Id          string    `json:"id" db:"id"`
	Description string    `json:"description,omitempty" db:"description"`
	Tag         string    `json:"tag,omitempty" db:"tag"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IImageTag interface {
	GetImageVersionById(id string) (ImageTag, error)
}
