package models

import (
	"time"

	db "github.com/sriharivishnu/shopify-challenge/services"
)

type ImageTag struct {
	Id           string    `json:"id" db:"id"`
	RepositoryId string    `json:"repository_id" db:"repository_id"`
	Description  string    `json:"description,omitempty" db:"description"`
	Tag          string    `json:"tag,omitempty" db:"tag"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IImageTag interface {
	Create() error
	GetImageTagById(id string) error
	GetLatestImageTag(id string) error
	GetImageTagsForRepo(repository_id string) ([]ImageTag, error)
}

func (tag *ImageTag) Create() error {
	tx := db.DbConn.MustBegin()
	err := tx.Get(tag, "INSERT INTO image_tag (repository_id, tag, description) VALUES (?, ?, ?) returning *;", tag.RepositoryId, tag.Tag, tag.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (tag *ImageTag) GetImageTagById(id string) error {
	err := db.DbConn.Get(tag, "select * from image_tag where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (tag *ImageTag) GetLatestImageTag(repository_id string) error {
	sql := "select * from image_tag where repository_id = ? order by created_at limit 1"
	err := db.DbConn.Get(tag, sql, repository_id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetImageTagsForRepo(repository_id string) ([]ImageTag, error) {
	imageTags := []ImageTag{}
	err := db.DbConn.Select(&imageTags, "select * from image_tag where repository_id = ?", repository_id)
	if err != nil {
		return nil, err
	}
	return imageTags, nil
}
