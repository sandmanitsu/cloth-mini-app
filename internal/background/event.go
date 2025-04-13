package background

import (
	domain "cloth-mini-app/internal/domain/event"
	ldomain "cloth-mini-app/internal/domain/lock"
	sl "cloth-mini-app/internal/logger"
	"context"
	"fmt"
	"log/slog"
	"time"
)

const (
	frequenceSendEvents = time.Second * 10
)

type EventBackground struct {
	logger     *slog.Logger
	outboxRepo OutboxRepository
	lockSrv    LockService
	producer   Producer
}

func NewEventBackground(logger *slog.Logger, otbxr OutboxRepository, ls LockService, prod Producer) *EventBackground {
	return &EventBackground{
		logger:     logger,
		outboxRepo: otbxr,
		lockSrv:    ls,
		producer:   prod,
	}
}

func (e *EventBackground) StartSendEvent() {
	const op = "background.event.StartSendEvent"
	e.logger.Info(fmt.Sprintf("%s: event send task started...", op))

	go func() {
		ticker := time.NewTicker(frequenceSendEvents)

		for {
			select {
			case <-ticker.C:
				ctx := context.Background()

				events, err := e.getEvents(ctx)
				if err != nil {
					e.logger.Error(fmt.Sprintf("%s : failed get events", op), sl.Err(err))
					continue
				}
				if len(events) == 0 {
					continue
				}

				successEventsId := make([]int, 0, len(events))
				for _, event := range events {
					if err := e.producer.WriteMesage(ctx, event.Payload); err != nil {
						e.logger.Error(fmt.Sprintf("%s : failed send event", op), sl.Err(err))
						continue
					}

					successEventsId = append(successEventsId, event.Id)
				}

				if len(successEventsId) != 0 {
					e.outboxRepo.ChangeStatus(ctx, successEventsId)
				}
			}
		}
	}()
}

func (e *EventBackground) getEvents(ctx context.Context) ([]domain.Event, error) {
	const op = "background.event.GetEvents"

	if err := e.lockSrv.AdvisoryLock(ctx, ldomain.OutboxAdvisoryLockId); err != nil {
		e.logger.Error(fmt.Sprintf("%s : failed get advisory lock", op), sl.Err(err))

		return nil, err
	}
	defer func() {
		if err := e.lockSrv.AdvisoryUnlock(ctx, ldomain.OutboxAdvisoryLockId); err != nil {
			e.logger.Error(fmt.Sprintf("%s : failed advisory unlock", op), sl.Err(err))
		}
	}()

	events, err := e.outboxRepo.GetEvents(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
