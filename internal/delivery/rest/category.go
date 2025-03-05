package rest

import (
	domain "cloth-mini-app/internal/domain/category"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CategoryService interface {
	GetCategories() ([]domain.Category, error)
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

type Category struct {
	CategoryId int    `json:"category_id"`
	Type       int    `json:"type"`
	Name       string `json:"category_name"`
}

func (c *CategoryHandler) Categories(ctx echo.Context) error {
	categories, err := c.Service.GetCategories()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "getting categories",
		})
	}

	categoriesResponse := make([]Category, 0, len(categories))
	for _, cat := range categories {
		categoriesResponse = append(categoriesResponse, Category{
			CategoryId: cat.CategoryId,
			Type:       cat.Type,
			Name:       cat.Name,
		})
	}

	return ctx.JSON(http.StatusOK, categoriesResponse)
}
