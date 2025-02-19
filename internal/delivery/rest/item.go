package rest

import (
	"cloth-mini-app/internal/domain"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type ItemService interface {
	// Fetching items
	Items(params url.Values) ([]domain.ItemAPI, error)
}

type ItemHandler struct {
	Service ItemService
}

type ErrorResponse struct {
	Err string `json:"error"`
}

func NewItemHandler(e *echo.Echo, srv ItemService) {
	handler := &ItemHandler{
		Service: srv,
	}

	g := e.Group("/item")
	g.GET("/get", handler.Items)
}

type ItemResponse struct {
	Count int              `json:"count"`
	Items []domain.ItemAPI `json:"items"`
}

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
