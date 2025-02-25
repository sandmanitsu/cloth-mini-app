package image

import (
	sl "cloth-mini-app/internal/logger"
	"log/slog"

	"github.com/google/uuid"
)

type MinioClient interface {
	CreateImage(file []byte, objectID string) error
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

// todo. Возвращать ссылку на изображение
func (i *ImageService) CreateItemImage(itemId int, file []byte) error {
	objectID := uuid.New().String()

	err := i.storage.CreateImage(file, objectID)
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
