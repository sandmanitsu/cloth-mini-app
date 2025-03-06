package image

import (
	"cloth-mini-app/internal/dto"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type MinioClient interface {
	Put(file dto.FileDTO) error
	Get(objectId string) (dto.FileDTO, error)
	GetMany(objectIds []string) ([]dto.FileDTO, error)
}

type ImageRepository interface {
	Insert(ctx context.Context, itemId int, objectID string) error
	Delete(ctx context.Context, imageId string) error
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
func (i *ImageService) CreateItemImage(ctx context.Context, itemId int, file []byte) (string, error) {
	objectID := uuid.New().String()

	err := i.storage.Put(dto.FileDTO{
		ID:          objectID,
		ContentType: minio.ImageContentType,
		Buffer:      file,
	})
	if err != nil {
		i.logger.Error("failet store image", sl.Err(err))
		return "", err
	}

	err = i.imageRepo.Insert(ctx, itemId, objectID)
	if err != nil {
		// todo. что тогда делать с изображением в хранилище s3???
		return "", err
	}

	return objectID, nil
}

// Get image from storage
func (i *ImageService) GetImage(ctx context.Context, imageId string) (file dto.FileDTO, err error) {
	file, err = i.storage.Get(imageId)
	if err != nil {
		i.logger.Error("failed getting image from storage", sl.Err(err))

		return
	}

	return
}

// Get images from storage
func (i *ImageService) GetImageMany(ctx context.Context, imageIds []string) ([]dto.FileDTO, error) {
	files, err := i.storage.GetMany(imageIds)
	if err != nil {
		i.logger.Error("failed getting image from storage", sl.Err(err))

		return nil, err
	}

	return files, nil
}

func (i *ImageService) Delete(ctx context.Context, imageId string) error {
	return i.imageRepo.Delete(ctx, imageId)
}
