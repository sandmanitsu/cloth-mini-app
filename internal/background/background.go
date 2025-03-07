package background

import (
	"cloth-mini-app/internal/storage/minio"
	"log/slog"
)

type BackgroundTask struct {
	TempImage *ImageBackground
}

func NewBackgroundTask(logger *slog.Logger, mc *minio.MinioClient, imr ImageRepository) *BackgroundTask {
	return &BackgroundTask{
		TempImage: NewImageBackground(logger, mc, imr),
	}
}
