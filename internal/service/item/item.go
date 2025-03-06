package item

import (
	domain "cloth-mini-app/internal/domain/item"
	sl "cloth-mini-app/internal/logger"
	"context"
	"fmt"
	"log/slog"
)

type ItemRepository interface {
	// Fetch items from db
	GetItems(ctx context.Context, params domain.ItemInputData) ([]domain.ItemAPI, error)
	// Returning item by id
	GetItemById(ctx context.Context, id int) (domain.ItemAPI, error)
	// Update item record
	Update(ctx context.Context, data domain.ItemUpdate) error
	// Create item
	Create(ctx context.Context, item domain.ItemCreate) error
	// Delete item
	Delete(ctx context.Context, id int) error
}

type ImageRepository interface {
	// Get images fileIds
	GetImages(ctx context.Context, itemId int) ([]string, error)
}

type ItemService struct {
	logger    *slog.Logger
	itemRepo  ItemRepository
	imageRepo ImageRepository
}

// Get item service object that represent the rest.ItemService interface
func NewItemService(logger *slog.Logger, ir ItemRepository, imr ImageRepository) *ItemService {
	return &ItemService{
		logger:    logger,
		itemRepo:  ir,
		imageRepo: imr,
	}
}

// Fetch items with provided params
func (i *ItemService) GetItems(ctx context.Context, params domain.ItemInputData) ([]domain.ItemAPI, error) {
	items, err := i.itemRepo.GetItems(ctx, params)
	if err != nil {
		return nil, err
	}

	return items, nil
}

type ItemUpdateData struct {
	ID   int
	Data map[string]any
}

// Prepare data to update
func (i *ItemService) Update(ctx context.Context, item domain.ItemUpdate) error {
	if item.ID == 0 {
		err := fmt.Errorf("error: empty or invalid id")
		i.logger.Error("update item", sl.Err(err))

		return err
	}

	err := i.itemRepo.Update(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (i *ItemService) GetItemById(ctx context.Context, id int) (domain.ItemAPI, error) {
	item, err := i.itemRepo.GetItemById(ctx, id)
	if err != nil {
		return item, err
	}

	imagesId, err := i.imageRepo.GetImages(ctx, int(item.ID))
	if err != nil {
		return item, err
	}

	item.ImageId = imagesId

	return item, err
}

func (i *ItemService) Create(ctx context.Context, item domain.ItemCreate) error {
	return i.itemRepo.Create(ctx, item)
}

func (i *ItemService) Delete(ctx context.Context, id int) error {
	return i.itemRepo.Delete(ctx, id)
}
