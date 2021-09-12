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
	FileKey      string    `json:"file_key,omitempty" db:"file_key"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IImageTag interface {
	Create() error
	GetImageTagByRepoAndTag(repository_id, tagName string) error
	GetLatestImageTag(id string) error
	GetImageTagsForRepo(repository_id string) ([]ImageTag, error)
}

func (tag *ImageTag) Create() error {
	sql := "INSERT INTO image_tag (repository_id, tag, description, file_key) VALUES (?, ?, ?, ?) returning *;"
	tx := db.DbConn.MustBegin()
	err := tx.Get(tag, sql, tag.RepositoryId, tag.Tag, tag.Description, tag.FileKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (tag *ImageTag) GetImageTagByRepoAndTag(repository_id, tagName string) error {
	err := db.DbConn.Get(tag, "select * from image_tag where repository_id = ? and tag = ?", repository_id, tagName)
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
