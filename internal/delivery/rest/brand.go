package rest

import (
	domain "cloth-mini-app/internal/domain/brand"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type BrandService interface {
	GetBrands() ([]domain.Brand, error)
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

type Brand struct {
	ID   int    `json:"brand_id"`
	Name string `json:"brand_name"`
}

func (b *BrandHandler) Brands(ctx echo.Context) error {
	brands, err := b.Service.GetBrands()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "getting brands",
		})
	}

	brandsResponse := make([]Brand, 0, len(brands))
	for _, brand := range brands {
		brandsResponse = append(brandsResponse, Brand{
			ID:   brand.ID,
			Name: brand.Name,
		})
	}

	return ctx.JSON(http.StatusOK, brandsResponse)
}
