package rest

import (
	"cloth-mini-app/internal/domain"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ItemService interface {
	// Fetching items
	Items(params url.Values) ([]domain.ItemAPI, error)
	Update(iten ItemDTO) error
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
	g.POST("/update/:id", handler.Update)
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
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "error: parse query params"})
	}

	items, err := i.Service.Items(request.Form)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "error: getting items"})
	}

	return c.JSON(http.StatusOK, ItemResponse{
		Count: len(items),
		Items: items,
	})
}

type ItemDTO struct {
	ID          int     `param:"id"`
	Brand       *string `json:"brand"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Sex         *int    `json:"sex"`
	CategoryId  *int    `json:"category_id"`
	Price       *int    `json:"price"`
	Discount    *int    `json:"discount"`
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
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "error: binding params"})
	}

	i.Service.Update(item)

	return c.JSON(http.StatusOK, ItemUpdateResponse{
		Success: true,
	})
}
