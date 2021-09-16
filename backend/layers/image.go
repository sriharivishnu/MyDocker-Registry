package layers

import (
	"github.com/sriharivishnu/shopify-challenge/models"
	db "github.com/sriharivishnu/shopify-challenge/services"
)

type ImageLayer interface {
	Create(repoId, tag, description, fileKey string) (models.ImageTag, error)
	GetImageTagByRepoAndTag(repository_id, tagName string) (models.ImageTag, error)
	GetImageTagsForRepo(repository_id string) ([]models.ImageTag, error)
}

type ImageService struct{}

func (service *ImageService) Create(repoId, tag, description, fileKey string) (models.ImageTag, error) {
	image := models.ImageTag{}
	sql := "INSERT INTO image_tag (repository_id, tag, description, file_key) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE description = VALUES(description), file_key = VALUES(file_key) returning *;"
	tx := db.DbConn.MustBegin()
	err := tx.Get(&image, sql, repoId, tag, description, fileKey)
	if err != nil {
		tx.Rollback()
		return image, err
	}
	tx.Commit()
	return image, nil
}

func (service *ImageService) GetImageTagByRepoAndTag(repository_id, tagName string) (models.ImageTag, error) {
	image := models.ImageTag{}
	err := db.DbConn.Get(&image, "select * from image_tag where repository_id = ? and tag = ?", repository_id, tagName)
	return image, err
}

func (service *ImageService) GetImageTagsForRepo(repository_id string) ([]models.ImageTag, error) {
	imageTags := []models.ImageTag{}
	err := db.DbConn.Select(&imageTags, "select * from image_tag where repository_id = ?", repository_id)
	if err != nil {
		return nil, err
	}
	return imageTags, nil
}
