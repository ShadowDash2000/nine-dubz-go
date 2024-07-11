package usecase

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type S3Interactor struct {
	S3Repository S3Repository
}

func (si *S3Interactor) GetS3Client() *s3.Client {
	return si.S3Repository.GetS3Client()
}

func (si *S3Interactor) PutObject(file *os.File, key string) (*s3.PutObjectOutput, error) {
	return si.S3Repository.PutObject(file, key)
}

func (si *S3Interactor) GetObject(key string, fileRange string) (*s3.GetObjectOutput, error) {
	return si.S3Repository.GetObject(key, fileRange)
}
