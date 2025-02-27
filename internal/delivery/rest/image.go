package rest

import (
	"cloth-mini-app/internal/dto"
	"io"
	"net/http"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ImageService interface {
	// Store image
	CreateItemImage(itemId int, file []byte) (string, error)
	// Get image from storage
	Image(imageId string) (dto.FileDTO, error)
}

type ImageHandler struct {
	Service ImageService
}

type CreateImageResponse struct {
	FileId string `json:"file_id"`
}

func NewImageHandler(e *echo.Echo, srv ImageService) {
	handler := &ImageHandler{
		Service: srv,
	}

	g := e.Group("/image")
	g.Use(middleware.Logger())

	g.POST("/create", handler.CreateItemImage)
	g.GET("/get/:image_id", handler.Image)
}

func (i *ImageHandler) CreateItemImage(c echo.Context) error {
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

	fileId, err := i.Service.CreateItemImage(itemId, imageBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failet store image. Maybe reached max image per item",
		})
	}

	return c.JSON(http.StatusOK, CreateImageResponse{
		FileId: fileId,
	})
}

type ImageId struct {
	Id string `param:"image_id"`
}

func (i *ImageHandler) Image(c echo.Context) error {
	var imageId ImageId
	err := c.Bind(&imageId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	file, err := i.Service.Image(imageId.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "getting image from storage",
		})
	}

	response := c.Response()

	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", file.ContentType)
	response.Write(file.Buffer)

	return nil
}
