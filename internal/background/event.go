package background

import (
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

				events, err := e.outboxRepo.GetEvents(ctx)
				if err != nil {
					e.logger.Error(fmt.Sprintf("%s : failed get events", op), sl.Err(err))
					continue
				}
				if len(events) == 0 {
					continue
				}

				successEventsId := make([]int, 0, len(events))
				for _, event := range events {
					if err := e.producer.WriteMesage(ctx, event.EventType, event.Payload); err != nil {
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
