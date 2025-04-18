package background

import (
	idomain "cloth-mini-app/internal/domain/image"
	ldomain "cloth-mini-app/internal/domain/lock"
	sl "cloth-mini-app/internal/logger"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const (
	frequenceCheckTempImages = time.Second * 10 // todo. поправить на time.Hour
	tempImageTTL             = time.Minute * 30 // todo. поправить на time.Hour
)

var (
	errNoImageToDelete = fmt.Errorf("no image to delete")
)

type ImageBackground struct {
	logger    *slog.Logger
	minioCl   *minio.MinioClient
	imageRepo ImageRepository
	lockSrv   LockService
}

func NewImageBackground(logger *slog.Logger, mc *minio.MinioClient, imr ImageRepository, lsrv LockService) *ImageBackground {
	return &ImageBackground{
		logger:    logger,
		minioCl:   mc,
		imageRepo: imr,
		lockSrv:   lsrv,
	}
}

func (i *ImageBackground) StartDeleteTempImage() {
	const op = "background.image.StartDeleteTempImage"
	i.logger.Info(fmt.Sprintf("%s: task started...", op))

	go func() {
		ticker := time.NewTicker(frequenceCheckTempImages)

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()

				if err := i.lockSrv.AdvisoryLock(ctx, ldomain.TempImageAdvisoryLockId); err != nil {
					i.logger.Error(fmt.Sprintf("%s : failed get advisory lock", op), sl.Err(err))
				}

				err := i.imageRepo.DeleteTempImage(ctx, func(images []idomain.TempImage) ([]idomain.TempImage, error) {
					if len(images) == 0 {
						return nil, errNoImageToDelete
					}

					deletingImage := make([]idomain.TempImage, 0, len(images))

					for _, image := range images {
						curr := time.Now()
						if curr.Sub(image.UploadedAt) > tempImageTTL {
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

				if err != nil && !errors.Is(err, errNoImageToDelete) {
					i.logger.Debug(op, sl.Err(err))
				}

				if err = i.lockSrv.AdvisoryUnlock(ctx, ldomain.TempImageAdvisoryLockId); err != nil {
					i.logger.Error(fmt.Sprintf("%s : failed advisory unlock", op), sl.Err(err))
				}
			}
		}
	}()
}
