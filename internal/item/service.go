package item

import (
	"cloth-mini-app/internal/domain"
	sl "cloth-mini-app/internal/logger"
	"log/slog"
	"net/url"
	"strconv"
)

type ItemRepository interface {
	// Fetch items from db
	Items(filter map[string]interface{}, limit, offset uint64) ([]domain.ItemAPI, error)
}

type ItemService struct {
	logger   *slog.Logger
	itemRepo ItemRepository
}

// Get service for item
func NewItemService(logger *slog.Logger, ir ItemRepository) *ItemService {
	return &ItemService{
		logger:   logger,
		itemRepo: ir,
	}
}

// Fetch items with provided params
func (i *ItemService) Items(params url.Values) ([]domain.ItemAPI, error) {
	var limitInt, offsetInt uint64
	var err error
	offset := params.Get("offset")
	limit := params.Get("limit")

	if limit != "" {
		limitInt, err = strconv.ParseUint(limit, 10, 64)
		if err != nil {
			i.logger.Error("error: limit to uint", sl.Err(err))
			return nil, err
		}
	}
	if offset != "" {
		offsetInt, err = strconv.ParseUint(offset, 10, 64)
		if err != nil {
			i.logger.Error("error: offset to uint", sl.Err(err))
			return nil, err
		}
	}

	filter := validateInputParams(params)

	items, err := i.itemRepo.Items(filter, limitInt, offsetInt)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Validating input params
func validateInputParams(params url.Values) map[string]interface{} {
	filter := make(map[string]interface{})

	if params.Get("id") != "" {
		filter["id"] = params.Get("id")
	}
	if params.Get("brand") != "" {
		filter["brand"] = params.Get("brand")
	}
	if params.Get("name") != "" {
		filter["name"] = params.Get("name")
	}
	if params.Get("sex") != "" {
		filter["sex"] = params.Get("sex")
	}
	if params.Get("category_type") != "" {
		filter["category_type"] = params.Get("category_type")
	}
	if params.Get("category_name") != "" {
		filter["category_name"] = params.Get("category_name")
	}
	if params.Get("price") != "" {
		filter["price"] = params.Get("price")
	}
	if params.Get("") != "" {
		filter["discount"] = params.Get("discount")
	}
	if params.Get("") != "" {
		filter["outer_link"] = params.Get("outer_link")
	}

	return filter
}
