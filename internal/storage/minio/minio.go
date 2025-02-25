package minio

import (
	"bytes"
	"cloth-mini-app/internal/config"
	"context"
	"fmt"
	"io"

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
func (m *MinioClient) CreateImage(file []byte, objectID string) error {
	const op = "storage.minio.CreateImage"

	reader := bytes.NewReader(file)
	_, err := m.cl.PutObject(context.Background(), m.bucketName, objectID, reader, int64(len(file)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *MinioClient) GetImage(objectId string) ([]byte, error) {
	const op = "storage.minio.GetImage"

	obj, err := m.cl.GetObject(context.Background(), m.bucketName, objectId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer obj.Close()

	objInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	buffer := make([]byte, objInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return buffer, nil
}
