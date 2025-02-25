package image

import (
	"cloth-mini-app/internal/dto"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type MinioClient interface {
	Put(file dto.FileDTO) error
	Get(objectId string) (dto.FileDTO, error)
}

type ImageRepository interface {
	Insert(itemId int, objectID string) error
}

type ImageService struct {
	logger    *slog.Logger
	storage   MinioClient
	imageRepo ImageRepository
}

func NewImageService(logger *slog.Logger, storage MinioClient, imageRepo ImageRepository) *ImageService {
	return &ImageService{
		logger:    logger,
		storage:   storage,
		imageRepo: imageRepo,
	}
}

// Put image to file storage storage and add file id to db
func (i *ImageService) CreateItemImage(itemId int, file []byte) error {
	objectID := uuid.New().String()

	err := i.storage.Put(dto.FileDTO{
		ID:          objectID,
		ContentType: minio.ImageContentType,
		Buffer:      file,
	})
	if err != nil {
		i.logger.Error("failet store image", sl.Err(err))
		return err
	}

	err = i.imageRepo.Insert(itemId, objectID)
	if err != nil {
		// todo. что тогда делать с изображением в хранилище s3???
		i.logger.Error("failet insert image to db", sl.Err(err))
		return err
	}

	return nil
}

// Get image from storage
func (i *ImageService) Image(imageId string) (file dto.FileDTO, err error) {
	file, err = i.storage.Get(imageId)
	if err != nil {
		i.logger.Error("failed getting image from storage", sl.Err(err))

		return
	}

	return
}

// Get many files from storage
func (i *ImageService) ImageMany(imageIds []string) ([]dto.FileDTO, error) {
	_, cancel := context.WithCancel(context.Background())

	workersCnt := 10
	var worker = func(imageIdCh <-chan string, fileCh chan<- dto.FileDTO, errCh chan<- error, wg *sync.WaitGroup) {
		for id := range imageIdCh {
			defer wg.Done()

			file, err := i.storage.Get(id)
			if err != nil {
				errCh <- err
				cancel()
				return
			}
			fileCh <- file

		}
	}

	var wg sync.WaitGroup
	imageIdCh := make(chan string, len(imageIds))
	fileCh := make(chan dto.FileDTO, len(imageIds))
	errCh := make(chan error, len(imageIds))

	for range workersCnt {
		go worker(imageIdCh, fileCh, errCh, &wg)
	}

	for _, id := range imageIds {
		wg.Add(1)
		imageIdCh <- id
	}
	close(imageIdCh)

	go func() {
		wg.Wait()
		close(fileCh)
		close(errCh)
	}()

	files := make([]dto.FileDTO, 0, len(imageIds))
	errors := make([]error, 0)
	for file := range fileCh {
		files = append(files, file)
	}
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		i.logger.Error("failed getting files", sl.Err(fmt.Errorf("%v", errors))) // todo. %v????
		return nil, fmt.Errorf("failed getting files")
	}

	return files, nil
}
