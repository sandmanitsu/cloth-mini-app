package item

import (
	domain "cloth-mini-app/internal/domain/item"
	sl "cloth-mini-app/internal/logger"
	"fmt"
	"log/slog"
)

type ItemRepository interface {
	// Fetch items from db
	GetItems(params domain.ItemInputData) ([]domain.ItemAPI, error)
	// Returning item by id
	GetItemById(id int) (domain.ItemAPI, error)
	// Update item record
	Update(data domain.ItemUpdate) error
	// Create item
	Create(item domain.ItemCreate) error
	// Delete item
	Delete(id int) error
}

type ImageRepository interface {
	// Get images fileIds
	GetImages(itemId int) ([]string, error)
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
func (i *ItemService) Items(params domain.ItemInputData) ([]domain.ItemAPI, error) {
	items, err := i.itemRepo.GetItems(params)
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
func (i *ItemService) Update(item domain.ItemUpdate) error {
	if item.ID == 0 {
		err := fmt.Errorf("error: empty or invalid id")
		i.logger.Error("update item", sl.Err(err))

		return err
	}

	err := i.itemRepo.Update(item)
	if err != nil {
		return err
	}

	return nil
}

func (i *ItemService) ItemById(id int) (domain.ItemAPI, error) {
	item, err := i.itemRepo.GetItemById(id)
	if err != nil {
		return item, err
	}

	imagesId, err := i.imageRepo.GetImages(int(item.ID))
	if err != nil {
		return item, err
	}

	item.ImageId = imagesId

	return item, err
}

func (i *ItemService) Create(item domain.ItemCreate) error {
	return i.itemRepo.Create(item)
}

func (i *ItemService) Delete(id int) error {
	return i.itemRepo.Delete(id)
}
