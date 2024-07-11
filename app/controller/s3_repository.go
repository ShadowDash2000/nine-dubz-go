package controller

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

type S3Repository struct {
	Bucket string
}

func (sr *S3Repository) GetS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "ru-1"
		o.BaseEndpoint = aws.String("https://s3.timeweb.cloud")
	})

	return client
}

func (sr *S3Repository) PutObject(file *os.File, key string) (*s3.PutObjectOutput, error) {
	client := sr.GetS3Client()

	output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return output, nil
}

func (sr *S3Repository) GetObject(key string, fileRange string) (*s3.GetObjectOutput, error) {
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
