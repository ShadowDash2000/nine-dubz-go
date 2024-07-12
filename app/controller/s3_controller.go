package controller

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"nine-dubz/app/model"
	"nine-dubz/app/usecase"
	"os"
)

type S3Controller struct {
	S3Interactor usecase.S3Interactor
}

func NewS3Controller() *S3Controller {
	accessKey, ok := os.LookupEnv("S3_ACCESS_KEY")
	if !ok {
		fmt.Println("S3_ACCESS_KEY environment variable not set")
	}
	secretKey, ok := os.LookupEnv("S3_SECRET_KEY")
	if !ok {
		fmt.Println("S3_SECRET_KEY environment variable not set")
	}
	bucket, ok := os.LookupEnv("S3_BUCKET")
	if !ok {
		fmt.Println("S3_BUCKET environment variable not set")
	}

	return &S3Controller{
		S3Interactor: usecase.S3Interactor{
			S3Repository: NewS3Repository(accessKey, secretKey, bucket),
		},
	}
}

func (sc *S3Controller) Read(key string, fileRange string) (*s3.GetObjectOutput, error) {
	return sc.S3Interactor.GetObject(key, fileRange)
}

func (sc *S3Controller) Upload(file *model.File) (*s3.PutObjectOutput, error) {
	reader, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return sc.S3Interactor.PutObject(reader, file.Name)
}
