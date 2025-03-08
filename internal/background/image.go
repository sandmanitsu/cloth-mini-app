package background

import (
	domain "cloth-mini-app/internal/domain/image"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/minio"
	"fmt"
	"log/slog"
	"time"
)

const (
	frequence = time.Second * 10 // todo. поправить на time.Hour
	ttl       = time.Minute * 30 // todo. поправить на time.Hour
)

type ImageRepository interface {
	GetTempImages() ([]domain.TempImage, error)
	DeleteTempImage(id uint) error
}

type ImageBackground struct {
	logger    *slog.Logger
	minioCl   *minio.MinioClient
	imageRepo ImageRepository
}

func NewImageBackground(logger *slog.Logger, mc *minio.MinioClient, imr ImageRepository) *ImageBackground {
	return &ImageBackground{
		logger:    logger,
		minioCl:   mc,
		imageRepo: imr,
	}
}

func (i *ImageBackground) StartDeleteTempImage() {
	const op = "background.image.StartDeleteTempImage"
	i.logger.Info(fmt.Sprintf("%s: task started...", op))

	go func() {
		ticker := time.NewTicker(frequence)

		for {
			select {
			case <-ticker.C:
				images, err := i.imageRepo.GetTempImages()
				if err != nil {
					continue
				}

				if len(images) == 0 {
					continue
				}

				for _, image := range images {
					curr := time.Now()
					if curr.Sub(image.UploadedAt) > ttl {
						err := i.imageRepo.DeleteTempImage(image.ID)
						if err != nil {
							i.logger.Error(fmt.Sprintf("%s: failed delete image from db", op), sl.Err(err))

							continue
						}

						err = i.minioCl.Delete(image.ObjectId)
						if err != nil {
							i.logger.Error(fmt.Sprintf("%s: failed delete image from s3", op), sl.Err(err))

							continue
						}
					}
				}
			}
		}
	}()
}
