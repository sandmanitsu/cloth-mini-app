package rest

import (
	"cloth-mini-app/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CategoryService interface {
	Categories() ([]domain.Category, error)
}

type CategoryHandler struct {
	Service CategoryService
}

func NewCategoryHandler(e *echo.Echo, srv CategoryService) {
	handler := &CategoryHandler{
		Service: srv,
	}

	g := e.Group("/category")
	g.Use(middleware.Logger())

	g.GET("/get", handler.Categories)
}

func (c *CategoryHandler) Categories(ctx echo.Context) error {
	categories, err := c.Service.Categories()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "getting categories",
		})
	}

	return ctx.JSON(http.StatusOK, categories)
}
