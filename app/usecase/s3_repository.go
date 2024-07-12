package usecase

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type S3Repository interface {
	GetS3Client() *s3.Client
	PutObject(file *os.File, key string) (*s3.PutObjectOutput, error)
	GetObject(key string, fileRange string) (*s3.GetObjectOutput, error)
}
