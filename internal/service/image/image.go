package image

import (
	sl "cloth-mini-app/internal/logger"
	"log/slog"
)

type MinioClient interface {
	CreateImage(file []byte) (string, error)
}

type ImageRepository interface {
	Insert(itemId int, url string) error
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
	url, err := i.storage.CreateImage(file)
	if err != nil {
		i.logger.Error("failet store image", sl.Err(err))
		return err
	}

	err = i.imageRepo.Insert(itemId, url)
	if err != nil {
		// todo. что тогда делать с изображением в хранилище s3???
		i.logger.Error("failet insert image to db", sl.Err(err))
		return err
	}

	return nil
}
