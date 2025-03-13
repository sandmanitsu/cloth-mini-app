package rest

import (
	domain "cloth-mini-app/internal/domain/item"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ItemService interface {
	// Fetching items
	GetItems(ctx context.Context, params domain.ItemInputData) ([]domain.ItemAPI, error)
	// Getting item by ID
	GetItemById(ctx context.Context, id int) (domain.ItemAPI, error)
	// Update item data
	Update(ctx context.Context, item domain.ItemUpdate) error
	// Create item
	Create(ctx context.Context, item domain.ItemCreate) error
	// Delete item
	Delete(ctx context.Context, id int) error
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

// GET /item/get Fetch items by query params
func (i *ItemHandler) Items(c echo.Context) error {
	var itemInput ItemQueryParams
	err := c.Bind(&itemInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	items, err := i.Service.GetItems(c.Request().Context(), domain.ItemInputData{
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

	return c.JSON(http.StatusOK, ItemsResponse{
		Count: len(items),
		Items: i.convertItemAPIFromDomain(items),
	})
}

func (i *ItemHandler) convertItemAPIFromDomain(domainItems []domain.ItemAPI) []ItemResponse {
	items := make([]ItemResponse, 0, len(domainItems))
	for _, item := range domainItems {
		items = append(items, ItemResponse{
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

// POST /item/update/:id Update item with provided id (required) and updating params
func (i *ItemHandler) Update(c echo.Context) error {
	var item ItemUpdate
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	i.Service.Update(c.Request().Context(), domain.ItemUpdate{
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

type ItemId struct {
	Id int `param:"id"`
}

func (i *ItemHandler) ItemById(c echo.Context) error {
	var itemId ItemId
	err := c.Bind(&itemId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	item, err := i.Service.GetItemById(c.Request().Context(), itemId.Id)
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

	err = i.Service.Create(c.Request().Context(), domain.ItemCreate{
		BrandId:     item.BrandId,
		Name:        item.Name,
		Description: item.Description,
		Sex:         item.Sex,
		CategoryId:  item.CategoryId,
		Price:       item.Price,
		Discount:    item.Discount,
		OuterLink:   item.OuterLink,
		Images:      item.Images,
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

	err = i.Service.Delete(c.Request().Context(), itemId.Id)
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
