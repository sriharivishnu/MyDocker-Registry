package external

import (
	"context"
	"log"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	localConfig "github.com/sriharivishnu/shopify-challenge/config"
	"github.com/sriharivishnu/shopify-challenge/utils"
)

type Storage interface {
	GetUploadURL(username string, repository string, tag string) (string, error)
	GetDownloadURL(file_key string) (string, error)
}

type S3 struct{}

type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)

	PresignPutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input)
}

func PutPresignedURL(c context.Context, api S3PresignGetObjectAPI, params *s3.PutObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignPutObject(c, params)
}

func (s *S3) GetUploadURL(username string, repository string, tag string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)
	file_key := utils.CreateFileKey(username, repository, tag)
	log.Println(file_key)

	input := &s3.PutObjectInput{
		Bucket: &localConfig.Config.S3_BUCKET_KEY,
		Key:    &file_key,
	}
	psClient := s3.NewPresignClient(client)
	resp, err := PutPresignedURL(context.TODO(), psClient, input)
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

func (s *S3) GetDownloadURL(file_key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: &localConfig.Config.S3_BUCKET_KEY,
		Key:    &file_key,
	}
	psClient := s3.NewPresignClient(client)

	resp, err := GetPresignedURL(context.TODO(), psClient, input)
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}
