package s3storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	AccessKey    string
	SecretKey    string
	Bucket       string
	Region       string
	BaseEndpoint string
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
	region, ok := os.LookupEnv("S3_REGION")
	if !ok {
		region = "ru-1"
	}
	baseEndpoint, ok := os.LookupEnv("S3_BASE_ENDPOINT")
	if !ok {
		baseEndpoint = *aws.String("https://s3.timeweb.cloud")
	}

	return &S3Storage{
		AccessKey:    accessKey,
		SecretKey:    secretKey,
		Bucket:       bucket,
		Region:       region,
		BaseEndpoint: baseEndpoint,
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
		o.Region = sr.Region
		o.BaseEndpoint = &sr.BaseEndpoint
	})

	return client
}

func (sr *S3Storage) MultipartUpload(ctx context.Context, file io.ReadSeeker, key, prefix string) (*s3.CompleteMultipartUploadOutput, int64, error) {
	client := sr.GetS3Client()
	key, err := url.JoinPath(prefix, key)
	if err != nil {
		return nil, 0, err
	}

	ctx, _ = context.WithCancel(ctx)
	createMultipartUploadOutput, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, 0, err
	}

	abortMultipartUploadInput := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(sr.Bucket),
		Key:      aws.String(key),
		UploadId: createMultipartUploadOutput.UploadId,
	}

	var completedParts []types.CompletedPart
	partNumber := int32(1)
	chunkSize := 1024 * 1024 * 50
	bytesRead := int64(0)
	buff := make([]byte, chunkSize)
	for {
		n, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}

			client.AbortMultipartUpload(ctx, abortMultipartUploadInput)
			return nil, 0, err
		}
		bytesRead += int64(n)

		uploadPartOutput, err := client.UploadPart(ctx, &s3.UploadPartInput{
			Bucket:        aws.String(sr.Bucket),
			Key:           aws.String(key),
			Body:          bytes.NewReader(buff[:n]),
			PartNumber:    aws.Int32(partNumber),
			UploadId:      createMultipartUploadOutput.UploadId,
			ContentLength: aws.Int64(int64(n)),
		})
		if err != nil {
			client.AbortMultipartUpload(ctx, abortMultipartUploadInput)
			return nil, 0, err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadPartOutput.ETag,
			PartNumber: aws.Int32(partNumber),
		})

		partNumber++
	}

	completeMultipartUploadOutput, err := client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(sr.Bucket),
		Key:      aws.String(key),
		UploadId: createMultipartUploadOutput.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		return nil, 0, err
	}

	return completeMultipartUploadOutput, bytesRead, nil
}

func (sr *S3Storage) PutObject(file io.ReadSeeker, key, prefix string) (*s3.PutObjectOutput, error) {
	client := sr.GetS3Client()
	key, err := url.JoinPath(prefix, key)
	if err != nil {
		return nil, err
	}

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

func (sr *S3Storage) GetRangeObject(key, prefix, fileRange string) (*s3.GetObjectOutput, error) {
	client := sr.GetS3Client()
	key, err := url.JoinPath(prefix, key)
	if err != nil {
		return nil, err
	}

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

func (sr *S3Storage) GetObject(key, prefix string) (*s3.GetObjectOutput, error) {
	client := sr.GetS3Client()
	key, err := url.JoinPath(prefix, key)
	if err != nil {
		return nil, err
	}

	object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (sr *S3Storage) DeleteAllInPrefix(prefix string) (*s3.DeleteObjectsOutput, error) {
	client := sr.GetS3Client()
	prefixPath, err := url.JoinPath(prefix, "/")
	if err != nil {
		return nil, err
	}

	directories := []string{prefixPath}
	var objects []types.ObjectIdentifier
	for len(directories) > 0 {
		for key, prefix := range directories {
			directories = slices.Delete(directories, key, 1)

			listObject, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
				Bucket: aws.String(sr.Bucket),
				Prefix: aws.String(prefix),
			})
			if err != nil {
				return nil, err
			}

			for _, object := range listObject.Contents {
				if strings.HasSuffix(*object.Key, "/") {
					directories = append(directories, *object.Key)
					continue
				}

				objects = append(objects, types.ObjectIdentifier{Key: object.Key})
			}
		}
	}

	deleteObjects, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(sr.Bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		return nil, err
	}

	return deleteObjects, nil
}

func (sr *S3Storage) DeleteObject(key, prefix string) (*s3.DeleteObjectOutput, error) {
	client := sr.GetS3Client()
	key, err := url.JoinPath(prefix, key)
	if err != nil {
		return nil, err
	}

	object, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(sr.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (sr *S3Storage) DeleteObjects(keys []string, prefix string) (*s3.DeleteObjectsOutput, error) {
	client := sr.GetS3Client()

	var objects []types.ObjectIdentifier
	for _, key := range keys {
		key, err := url.JoinPath(prefix, key)
		if err != nil {
			return nil, err
		}

		objects = append(objects, types.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	object, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(sr.Bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}
