package rest

import (
	"cloth-mini-app/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type BrandService interface {
	Brands() ([]domain.Brand, error)
}

type BrandHandler struct {
	Service BrandService
}

func NewBrandHandler(e *echo.Echo, srv BrandService) {
	handler := &BrandHandler{
		Service: srv,
	}

	g := e.Group("/brand")
	g.Use(middleware.Logger())

	g.GET("/get", handler.Brands)
}

func (b *BrandHandler) Brands(ctx echo.Context) error {
	categories, err := b.Service.Brands()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "getting categories",
		})
	}

	return ctx.JSON(http.StatusOK, categories)
}
