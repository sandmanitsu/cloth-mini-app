package rest

import (
	"cloth-mini-app/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ItemService interface {
	// Fetching items
	Items(params url.Values) ([]domain.ItemAPI, error)
	// Getting item by ID
	ItemById(id int) (domain.ItemAPI, error)
	// Update item data
	Update(iten ItemDTO) error
	// Create item
	Create(item ItemCreateDTO) error
}

type ItemHandler struct {
	Service ItemService
}

type ErrorResponse struct {
	Err string `json:"error"`
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
}

type ItemResponse struct {
	Count int              `json:"count"`
	Items []domain.ItemAPI `json:"items"`
}

// GET /item/get Fetch items by query params
func (i *ItemHandler) Items(c echo.Context) error {
	request := c.Request()
	err := request.ParseForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "parse query params"})
	}

	items, err := i.Service.Items(request.Form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "getting items"})
	}

	return c.JSON(http.StatusOK, ItemResponse{
		Count: len(items),
		Items: items,
	})
}

type ItemDTO struct {
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

type ItemUpdateResponse struct {
	Success bool `json:"update"`
}

// POST /item/update/:id Update item with provided id (required) and updating params
func (i *ItemHandler) Update(c echo.Context) error {
	var item ItemDTO
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	i.Service.Update(item)

	return c.JSON(http.StatusOK, ItemUpdateResponse{
		Success: true,
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

	item, err := i.Service.ItemById(itemId.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "no records with provided id"})
		}
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "failed getting item"})
	}

	return c.JSON(http.StatusOK, item)
}

type ItemCreateDTO struct {
	BrandId     int    `json:"brand_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sex         int    `json:"sex" validate:"required"`
	CategoryId  int    `json:"category_id" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	Discount    uint   `json:"discount"`
	OuterLink   string `json:"outer_link" validate:"required"`
}

type ItemCreateResponse struct {
	Success bool `json:"create"`
}

func (i *ItemHandler) Create(c echo.Context) error {
	var item ItemCreateDTO
	err := c.Bind(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(item); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: fmt.Sprintf("validation params : %s", err)})
	}

	err = i.Service.Create(item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "failed creating item"})
	}

	return c.JSON(http.StatusOK, ItemCreateResponse{
		Success: true,
	})
}
