package minio

import (
	"bytes"
	"cloth-mini-app/internal/config"
	"cloth-mini-app/internal/dto"
	"context"
	"fmt"
	"io"
	"sync"

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

// Get many files from storage
// Use worker pull with 10 workers (amount workers - workersCnt)
func (m *MinioClient) GetMany(objectIds []string) ([]dto.FileDTO, error) {
	_, cancel := context.WithCancel(context.Background())

	workersCnt := 10
	var worker = func(objectIdCh <-chan string, fileCh chan<- dto.FileDTO, errCh chan<- error, wg *sync.WaitGroup) {
		for id := range objectIdCh {
			defer wg.Done()

			file, err := m.Get(id)
			if err != nil {
				errCh <- err
				cancel()
				return
			}
			fileCh <- file

		}
	}

	var wg sync.WaitGroup
	objectIdCh := make(chan string, len(objectIds))
	fileCh := make(chan dto.FileDTO, len(objectIds))
	errCh := make(chan error, len(objectIds))

	for range workersCnt {
		go worker(objectIdCh, fileCh, errCh, &wg)
	}

	for _, id := range objectIds {
		wg.Add(1)
		objectIdCh <- id
	}
	close(objectIdCh)

	go func() {
		wg.Wait()
		close(fileCh)
		close(errCh)
	}()

	files := make([]dto.FileDTO, 0, len(objectIds))
	errors := make([]error, 0)
	for file := range fileCh {
		files = append(files, file)
	}
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed getting files: err list: %v", errors)
	}

	return files, nil
}
