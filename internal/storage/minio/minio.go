package minio

import (
	"bytes"
	"cloth-mini-app/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	bucketName string
	cl         *minio.Client
}

// Create minio client object
func NewMinioClient(cfg config.Minio) (*MinioClient, error) {
	const op = "storage.minio.New"

	ctx := context.Background()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exist, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exist {
		err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &MinioClient{
		cl:         client,
		bucketName: cfg.BucketName,
	}, nil
}

// Put image to store and return url to the image
func (m *MinioClient) Create(file []byte) (string, error) {
	const op = "storage.minio.Create"

	objectID := uuid.New().String()

	reader := bytes.NewReader(file)
	_, err := m.cl.PutObject(context.Background(), m.bucketName, objectID, reader, int64(len(file)), minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	url, err := m.cl.PresignedGetObject(context.Background(), m.bucketName, objectID, time.Second*60*60*24, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url.String(), nil
}
