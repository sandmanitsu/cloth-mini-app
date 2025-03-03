package rest

import (
	domain "cloth-mini-app/internal/domain/item"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ItemService interface {
	// Fetching items
	Items(params domain.ItemInputData) ([]domain.ItemAPI, error)
	// Getting item by ID
	ItemById(id int) (domain.ItemAPI, error)
	// Update item data
	Update(item domain.ItemUpdate) error
	// Create item
	Create(item domain.ItemCreate) error
	// Delete item
	Delete(id int) error
}

type ItemHandler struct {
	Service ItemService
}

// Create item handler object
func NewItemHandler(e *echo.Echo, srv ItemService) {
	handler := &ItemHandler{
		Service: srv,
	}

	g := e.Group("/item")
	g.Use(middleware.Logger())
	g.GET("/get", handler.Items)
	g.GET("/get/:id", handler.ItemById)
	g.POST("/update/:id", handler.Update)
	g.POST("/create", handler.Create)
	g.DELETE("/delete/:id", handler.Delete)
}

type ItemInputParams struct {
	ID         *uint   `query:"id"`
	BrandId    *uint   `query:"brand_id"`
	Name       *string `query:"name"`
	Sex        *int    `query:"sex"`
	CategoryId *uint   `query:"category_id"`
	MinPrice   *uint   `query:"min_price"`
	MaxPrice   *uint   `query:"max_price"`
	Discount   *uint   `query:"discount"`
	Offset     *uint   `query:"offset"`
	Limit      *uint   `query:"limit"`
}

type Item struct {
	ID           uint       `json:"id"`
	BrandId      uint       `json:"brand_id"`
	BrandName    string     `json:"brand_name"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Sex          int        `json:"sex"`
	CategoryId   int        `json:"category_id"`
	CategoryType int        `json:"category_type"`
	CategoryName string     `json:"category_name"`
	Price        int        `json:"price"`
	Discount     *int       `json:"discount"`
	OuterLink    string     `json:"outer_link"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type ItemResponse struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

// GET /item/get Fetch items by query params
func (i *ItemHandler) Items(c echo.Context) error {
	var itemInput ItemInputParams
	err := c.Bind(&itemInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	items, err := i.Service.Items(domain.ItemInputData{
		ID:         itemInput.ID,
		BrandId:    itemInput.BrandId,
		Name:       itemInput.Name,
		Sex:        itemInput.Sex,
		CategoryId: itemInput.CategoryId,
		MinPrice:   itemInput.MinPrice,
		MaxPrice:   itemInput.MaxPrice,
		Discount:   itemInput.Discount,
		Offset:     itemInput.Offset,
		Limit:      itemInput.Limit,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "getting items"})
	}

	return c.JSON(http.StatusOK, ItemResponse{
		Count: len(items),
		Items: i.convertDomainToDeliveryItem(items),
	})
}

func (i *ItemHandler) convertDomainToDeliveryItem(domainItems []domain.ItemAPI) []Item {
	items := make([]Item, 0, len(domainItems))
	for _, item := range domainItems {
		items = append(items, Item{
			ID:           item.ID,
			BrandId:      item.BrandId,
			BrandName:    item.BrandName,
			Name:         item.Name,
			Description:  item.Description,
			Sex:          item.Sex,
			CategoryId:   item.CategoryId,
			CategoryType: item.CategoryType,
			CategoryName: item.CategoryName,
			Price:        item.Price,
			Discount:     item.Discount,
			OuterLink:    item.OuterLink,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return items
}

type ItemUpdate struct {
	ID          int     `param:"id"`
	BrandId     *int    `json:"brand_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Sex         *int    `json:"sex"`
	CategoryId  *int    `json:"category_id"`
	Price       *uint   `json:"price"`
	Discount    *uint   `json:"discount"`
	OuterLink   *string `json:"outerlink"`
}

// POST /item/update/:id Update item with provided id (required) and updating params
func (i *ItemHandler) Update(c echo.Context) error {
	var item ItemUpdate
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	i.Service.Update(domain.ItemUpdate{
		ID:          item.ID,
		BrandId:     item.BrandId,
		Name:        item.Name,
		Description: item.Description,
		Sex:         item.Sex,
		CategoryId:  item.CategoryId,
		Price:       item.Price,
		Discount:    item.Discount,
		OuterLink:   item.OuterLink,
	})

	return c.JSON(http.StatusOK, SuccessResponse{
		Status:    true,
		Operation: "update",
	})
}

type ItemByIdResponse struct {
	ID           uint       `json:"id"`
	BrandId      uint       `json:"brand_id"`
	BrandName    string     `json:"brand_name"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Sex          int        `json:"sex"`
	CategoryId   int        `json:"category_id"`
	CategoryType int        `json:"category_type"`
	CategoryName string     `json:"category_name"`
	Price        int        `json:"price"`
	Discount     *int       `json:"discount"`
	OuterLink    string     `json:"outer_link"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ImageId      []string   `json:"image_id"`
}

type ItemId struct {
	Id int `param:"id"`
}

func (i *ItemHandler) ItemById(c echo.Context) error {
	var itemId ItemId
	err := c.Bind(&itemId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	item, err := i.Service.ItemById(itemId.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "no records with provided id"})
		}
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "failed getting item"})
	}

	return c.JSON(http.StatusOK, ItemByIdResponse{
		ID:           item.ID,
		BrandId:      item.BrandId,
		BrandName:    item.BrandName,
		Name:         item.Name,
		Description:  item.Description,
		Sex:          item.Sex,
		CategoryId:   item.CategoryId,
		CategoryType: item.CategoryType,
		CategoryName: item.CategoryName,
		Price:        item.Price,
		Discount:     item.Discount,
		OuterLink:    item.OuterLink,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
		ImageId:      item.ImageId,
	})
}

type ItemCreate struct {
	BrandId     int    `json:"brand_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sex         int    `json:"sex" validate:"required"`
	CategoryId  int    `json:"category_id" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	Discount    uint   `json:"discount"`
	OuterLink   string `json:"outer_link" validate:"required"`
}

func (i *ItemHandler) Create(c echo.Context) error {
	var item ItemCreate
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(item); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: fmt.Sprintf("validation params : %s", err)})
	}

	err = i.Service.Create(domain.ItemCreate{
		BrandId:     item.BrandId,
		Name:        item.Name,
		Description: item.Description,
		Sex:         item.Sex,
		CategoryId:  item.CategoryId,
		Price:       item.Price,
		Discount:    item.Discount,
		OuterLink:   item.OuterLink,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "failed creating item"})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Status:    true,
		Operation: "create",
	})
}

func (i *ItemHandler) Delete(c echo.Context) error {
	var itemId ItemId
	err := c.Bind(&itemId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	err = i.Service.Delete(itemId.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failed deleting item",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Status:    true,
		Operation: "delete",
	})
}
