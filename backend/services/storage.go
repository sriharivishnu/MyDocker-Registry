package services

type Storage interface {
	GetUploadURL(repository string, imageName string, tag string)
	GetDownloadURL(repository string, imageName string, tag string)
}
