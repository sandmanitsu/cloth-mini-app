package background

import (
	domain "cloth-mini-app/internal/domain/image"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"fmt"
	"log/slog"
	"time"
)

const (
	frequence = time.Second * 10 // todo. поправить на time.Hour
	ttl       = time.Minute * 30 // todo. поправить на time.Hour
)

var (
	errNoImageToDelete = fmt.Errorf("no image to delete")
)

type ImageRepository interface {
	// Delete temp images data into db
	DeleteTempImage(ctx context.Context, deleteFn func([]domain.TempImage) ([]domain.TempImage, error)) error
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
				ctx := context.Background()

				i.imageRepo.DeleteTempImage(ctx, func(images []domain.TempImage) ([]domain.TempImage, error) {
					if len(images) == 0 {
						return nil, errNoImageToDelete
					}

					deletingImage := make([]domain.TempImage, 0, len(images))

					for _, image := range images {
						curr := time.Now()
						if curr.Sub(image.UploadedAt) > ttl {
							deletingImage = append(deletingImage, image)
							err := i.minioCl.Delete(ctx, image.ObjectId)
							if err != nil {
								i.logger.Error(fmt.Sprintf("%s: failed delete image from s3", op), sl.Err(err))

								return nil, err
							}
						}
					}

					return deletingImage, nil
				})
			}
		}
	}()
}
