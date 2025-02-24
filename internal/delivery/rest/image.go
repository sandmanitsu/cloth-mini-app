package rest

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ImageService interface {
	Create(itemId int, file []byte) error
}

type ImageHandler struct {
	Service ImageService
}

func NewImageHandler(e *echo.Echo, srv ImageService) {
	handler := &ImageHandler{
		Service: srv,
	}

	g := e.Group("/image")
	g.Use(middleware.Logger())

	g.POST("/create", handler.Create)
}

func (i *ImageHandler) Create(c echo.Context) error {
	request := c.Request()
	err := request.ParseForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "parse query params"})
	}

	itemId, err := strconv.Atoi(request.Form.Get("itemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "itemId is incorrect or not provided"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failed get file",
		})
	}

	image, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failed open file",
		})
	}
	defer image.Close()

	imageBytes, err := io.ReadAll(image)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failed read file",
		})
	}

	mtype := mimetype.Detect(imageBytes)
	if !(mtype.Is("image/jpeg") || mtype.Is("image/png")) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "incorrect image format. allowed image formats: .jpg/.png",
		})
	}

	if err = i.Service.Create(itemId, imageBytes); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failet store image",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Status:    true,
		Operation: "create",
	})
}
