package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewClient(endpoint string, bucketName string, accessKeyID string, secretAccessKey string, useSSL bool) (S3Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return S3Client{}, err
	}

	var client = S3Client{
		MinioClient: minioClient,
		BucketName:  bucketName,
	}

	return client, nil
}

func Download(ctx context.Context, client S3Client, content []byte, slug string) error {
	log.Println(client.BucketName)

	// TODO: stream this as response

	obj, err := client.MinioClient.GetObject(ctx, client.BucketName, fmt.Sprintf("%s.png", slug), minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	// TODO: pull some telemetry from info var

	log.Printf("Successfully uploaded %s of size %d\n", slug, info.Size)
	return nil
}
