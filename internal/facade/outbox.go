package facade

import (
	bdomain "cloth-mini-app/internal/domain/brand"
	edomain "cloth-mini-app/internal/domain/event"
	idomain "cloth-mini-app/internal/domain/item"
	"cloth-mini-app/internal/storage/postgresql"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event edomain.Event) error
}

type ItemImageRepository interface {
	Create(ctx context.Context, item idomain.ItemCreate) (uint, error)
}

type BrandRepositury interface {
	GetBrand(ctx context.Context, brandId int) (bdomain.Brand, error)
}

type OutboxFacade struct {
	db            *sql.DB
	logger        *slog.Logger
	outboxRepo    OutboxRepository
	itemImageRepo ItemImageRepository
	brandRepo     BrandRepositury
}

func NewOutboxFacade(db *postgresql.Storage, logger *slog.Logger, outboxr OutboxRepository, itimr ItemImageRepository, br BrandRepositury) *OutboxFacade {
	return &OutboxFacade{
		db:            db.DB,
		logger:        logger,
		outboxRepo:    outboxr,
		itemImageRepo: itimr,
		brandRepo:     br,
	}
}

func (o *OutboxFacade) CreateItemWithNotification(ctx context.Context, item idomain.ItemCreate) error {
	brand, err := o.brandRepo.GetBrand(ctx, item.BrandId)
	if err != nil {
		return err
	}

	err = postgresql.WrapTx(ctx, o.db, func(ctx context.Context) error {
		itemId, err := o.itemImageRepo.Create(ctx, item)
		if err != nil {
			return err
		}

		event := edomain.Event{
			EventType: edomain.EventCreateItem,
			Payload: []byte(fmt.Sprintf(
				`{"item_id":"%d", "brand_name":"%s", "item_name":"%s", "price":"%d"}`,
				itemId,
				brand.Name,
				item.Name,
				item.Price,
			)),
		}

		err = o.outboxRepo.CreateEvent(ctx, event)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
