package item

import "log/slog"

type ItemRepository interface {
	Items()
}

type ItemService struct {
	logger   *slog.Logger
	itemRepo ItemRepository
}

func NewItemService(logger *slog.Logger, ir ItemRepository) *ItemService {
	return &ItemService{
		logger:   logger,
		itemRepo: ir,
	}
}

func (i *ItemService) Items() {
	i.itemRepo.Items()
}
