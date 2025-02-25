package minio

import (
	"bytes"
	"cloth-mini-app/internal/config"
	"cloth-mini-app/internal/dto"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	ImageContentType = "image/jpeg"
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

// Put file to store
func (m *MinioClient) Put(file dto.FileDTO) error {
	const op = "storage.minio.CreateImage"

	reader := bytes.NewReader(file.Buffer)
	_, err := m.cl.PutObject(context.Background(), m.bucketName, file.ID, reader, int64(len(file.Buffer)), minio.PutObjectOptions{
		ContentType: file.ContentType,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Get file from storage
func (m *MinioClient) Get(objectId string) (dto.FileDTO, error) {
	const op = "storage.minio.GetImage"

	obj, err := m.cl.GetObject(context.Background(), m.bucketName, objectId, minio.GetObjectOptions{})
	if err != nil {
		return dto.FileDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	defer obj.Close()

	objInfo, err := obj.Stat()
	if err != nil {
		return dto.FileDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	buffer := make([]byte, objInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return dto.FileDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.FileDTO{
		ID:          objectId,
		ContentType: objInfo.ContentType,
		Buffer:      buffer,
	}, nil
}
