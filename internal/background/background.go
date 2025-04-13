package background

import (
	edomain "cloth-mini-app/internal/domain/event"
	idomain "cloth-mini-app/internal/domain/image"
	ldomain "cloth-mini-app/internal/domain/lock"
	"cloth-mini-app/internal/storage/minio"
	"context"
	"log/slog"
)

type BackgroundTask struct {
	TempImage *ImageBackground
	Event     *EventBackground
}

type ImageRepository interface {
	// Delete temp images data into db
	DeleteTempImage(ctx context.Context, deleteFn func([]idomain.TempImage) ([]idomain.TempImage, error)) error
}

type LockService interface {
	AdvisoryLock(ctx context.Context, id ldomain.AdvisoryLockId) error
	AdvisoryUnlock(ctx context.Context, id ldomain.AdvisoryLockId) error
}

type OutboxRepository interface {
	GetEvents(ctx context.Context) ([]edomain.Event, error)
	ChangeStatus(ctx context.Context, eventsId []int) error
}

type Producer interface {
	WriteMesage(ctx context.Context, payload []byte) error
}

func NewBackgroundTask(
	logger *slog.Logger,
	mc *minio.MinioClient,
	imr ImageRepository,
	lcrv LockService,
	outboxr OutboxRepository,
	producer Producer,
) *BackgroundTask {
	return &BackgroundTask{
		TempImage: NewImageBackground(logger, mc, imr, lcrv),
		Event:     NewEventBackground(logger, outboxr, lcrv, producer),
	}
}
