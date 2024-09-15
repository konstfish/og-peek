package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var s3Client *S3Client

type S3Client struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewClient(endpoint string, bucketName string, accessKeyID string, secretAccessKey string, useSSL bool) error {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	s3Client = &S3Client{
		MinioClient: minioClient,
		BucketName:  bucketName,
	}

	return nil
}

func Download(ctx context.Context, domain string, slug string) (*minio.Object, error) {
	obj, err := s3Client.MinioClient.GetObject(ctx, s3Client.BucketName, fmt.Sprintf("%s/%s.png", domain, slug), minio.GetObjectOptions{})
	if err != nil {
		log.Println("error downloading", slug, err)
		return nil, err
	}

	return obj, nil
}
