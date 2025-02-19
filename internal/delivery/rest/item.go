package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ItemService interface {
	Items()
}

type ItemHandler struct {
	Service ItemService
}

func NewItemHandler(e *echo.Echo, srv ItemService) {
	handler := &ItemHandler{
		Service: srv,
	}

	g := e.Group("/item")
	g.GET("/get", handler.Items)
}

func (i *ItemHandler) Items(c echo.Context) error {
	i.Service.Items()

	return c.JSON(http.StatusOK, "{\"status\":1}")
}
