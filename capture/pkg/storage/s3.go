package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

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

func Upload(ctx context.Context, content []byte, domain string, slug string) error {
	contentType := "image/png"

	// store image in bucket, set expiry to
	info, err := s3Client.MinioClient.PutObject(
		ctx,
		s3Client.BucketName,
		fmt.Sprintf("%s/%s.png", domain, slug),
		bytes.NewReader(content),
		int64(len(content)),
		minio.PutObjectOptions{
			ContentType: contentType,
			Expires:     time.Now().Add(7 * 24 * time.Hour),
		},
	)
	if err != nil {
		return err
	}

	// TODO: pull some telemetry from info var

	log.Printf("Successfully uploaded %s of size %d\n", slug, info.Size)
	return nil
}
