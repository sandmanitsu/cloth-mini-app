package item

import (
	"cloth-mini-app/internal/delivery/rest"
	"cloth-mini-app/internal/domain"
	sl "cloth-mini-app/internal/logger"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
)

type ItemRepository interface {
	// Fetch items from db
	Items(filter map[string]interface{}, limit, offset uint64) ([]domain.ItemAPI, error)
	// Update item record
	Update(data ItemUpdateData) error
	// Returning item by id
	ItemById(id int) (domain.ItemAPI, error)
	// Create item
	Create(item rest.ItemCreateDTO) error
}

type ItemService struct {
	logger   *slog.Logger
	itemRepo ItemRepository
}

// Get item service object that represent the rest.ItemService interface
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

	filter := i.validateInputParams(params)

	items, err := i.itemRepo.Items(filter, limitInt, offsetInt)
	if err != nil {
		i.logger.Debug("", sl.Err(err))
		return nil, err
	}

	return items, nil
}

// Validating input params
func (i *ItemService) validateInputParams(params url.Values) map[string]any {
	filter := make(map[string]any)

	if params.Get("id") != "" {
		filter["i.id"] = params.Get("id")
	}
	if params.Get("brand_id") != "" {
		filter["i.brand_id"] = params.Get("brand_id")
	}
	if params.Get("name") != "" {
		filter["i.name"] = params.Get("name")
	}
	if params.Get("sex") != "" {
		filter["i.sex"] = params.Get("sex")
	}
	if params.Get("category_id") != "" {
		filter["c.id"] = params.Get("category_id")
	}
	if params.Get("category_type") != "" {
		filter["c.category_type"] = params.Get("category_type")
	}
	if params.Get("category_name") != "" {
		filter["c.category_name"] = params.Get("category_name")
	}

	if params.Get("min_price") != "" || params.Get("max_price") != "" {
		fmt.Println(params.Get("min_price"), params.Get("max_price"))
		filter["min_price"] = params.Get("min_price")
		filter["max_price"] = params.Get("max_price")
	} else if params.Get("price") != "" {
		filter["i.price"] = params.Get("price")
	}

	if params.Get("") != "" {
		filter["i.discount"] = params.Get("discount")
	}
	if params.Get("") != "" {
		filter["i.outer_link"] = params.Get("outer_link")
	}

	return filter
}

type ItemUpdateData struct {
	ID   int
	Data map[string]any
}

// Prepare data to update
func (i *ItemService) Update(item rest.ItemDTO) error {
	if item.ID == 0 {
		return fmt.Errorf("error: empty or invalid id")
	}

	data := i.validateUpdateData(item)
	if len(data) == 0 {
		return nil
	}

	err := i.itemRepo.Update(ItemUpdateData{
		ID:   item.ID,
		Data: data,
	})
	if err != nil {
		return err
	}

	return nil
}

// Validating input update params
func (i *ItemService) validateUpdateData(item rest.ItemDTO) map[string]any {
	data := make(map[string]any)

	if item.BrandId != nil {
		data["brand_id"] = *item.BrandId
	}
	if item.Name != nil {
		data["name"] = *item.Name
	}
	if item.Description != nil {
		data["description"] = *item.Description
	}
	if item.CategoryId != nil {
		data["category_id"] = *item.CategoryId
	}
	if item.Sex != nil {
		data["sex"] = *item.Sex
	}
	if item.Discount != nil {
		data["discount"] = *item.Discount
	}
	if item.Price != nil {
		data["price"] = *item.Price
	}
	if item.OuterLink != nil {
		data["outer_link"] = *item.OuterLink
	}

	return data
}

func (i *ItemService) ItemById(id int) (domain.ItemAPI, error) {
	return i.itemRepo.ItemById(id)
}

func (i *ItemService) Create(item rest.ItemCreateDTO) error {
	return i.itemRepo.Create(item)
}
