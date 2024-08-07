package s3storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
)

type S3Storage struct {
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewS3Storage() *S3Storage {
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

	return &S3Storage{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
	}
}

func (sr *S3Storage) GetS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(sr.AccessKey, sr.SecretKey, ""),
		),
	)
	if err != nil {
		fmt.Println("unable to load SDK config, %v", err)
		return nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "ru-1"
		o.BaseEndpoint = aws.String("https://s3.timeweb.cloud")
	})

	return client
}

func (sr *S3Storage) PutObject(file io.ReadSeeker, key string) (*s3.PutObjectOutput, error) {
	client := sr.GetS3Client()

	output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (sr *S3Storage) GetRangeObject(key string, fileRange string) (*s3.GetObjectOutput, error) {
	client := sr.GetS3Client()

	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
		Range:  aws.String(fileRange),
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (sr *S3Storage) GetObject(key string) (*s3.GetObjectOutput, error) {
	client := sr.GetS3Client()

	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (sr *S3Storage) DeleteObject(key string) (*s3.DeleteObjectOutput, error) {
	client := sr.GetS3Client()

	object, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}
