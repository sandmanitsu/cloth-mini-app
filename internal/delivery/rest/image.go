package rest

import (
	"cloth-mini-app/internal/dto"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	errGetFile   = fmt.Errorf("failed get file")
	errOpenFile  = fmt.Errorf("failed open file")
	errReadFile  = fmt.Errorf("failed read file")
	errImageType = fmt.Errorf("incorrect image format. allowed image formats: .jpg/.png")
)

type ImageService interface {
	// Store image
	CreateItemImage(ctx context.Context, itemId int, file []byte) (string, error)
	// Get image from storage
	GetImage(ctx context.Context, imageId string) (dto.FileDTO, error)
	// Delete image from db and storage
	Delete(ctx context.Context, imageId string) error
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

	g.POST("/create", handler.CreateItemImage)
	g.POST("/temp", handler.CreateTempImage)
	g.GET("/get/:image_id", handler.Image)
	g.DELETE("/delete", handler.Delete)
}

type CreateImageResponse struct {
	FileId string `json:"file_id"`
}

// Put image to storage and adding fileID to db
// Return CreateImageResponse
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

	imageBytes, err := i.file(c)
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

	fileId, err := i.Service.CreateItemImage(c.Request().Context(), itemId, imageBytes)
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

// Return image by image_id in query param
func (i *ImageHandler) Image(c echo.Context) error {
	var imageId ImageId
	err := c.Bind(&imageId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "binding params"})
	}

	file, err := i.Service.GetImage(c.Request().Context(), imageId.Id)
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

// Deleting image by image_id provided in query param
func (i *ImageHandler) Delete(c echo.Context) error {
	request := c.Request()
	err := request.ParseForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "parse query params"})
	}

	imageId := request.Form.Get("image_id")
	if imageId == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "image_id not provided"})
	}

	err = i.Service.Delete(c.Request().Context(), imageId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Err: "failed deleting image"})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Status:    true,
		Operation: "delete",
	})
}

func (i *ImageHandler) CreateTempImage(c echo.Context) error {
	imageBytes, err := i.file(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: err.Error(),
		})
	}

	fileId, err := i.Service.CreateTempImage(imageBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Err: "failet store image. Maybe reached max image per item",
		})
	}

	return c.JSON(http.StatusOK, CreateImageResponse{
		FileId: fileId,
	})
}

// read image file
func (i *ImageHandler) file(c echo.Context) ([]byte, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return nil, errGetFile
	}

	image, err := file.Open()
	if err != nil {
		return nil, errOpenFile
	}
	defer image.Close()

	imageBytes, err := io.ReadAll(image)
	if err != nil {
		return nil, errReadFile
	}

	mtype := mimetype.Detect(imageBytes)
	if !(mtype.Is("image/jpeg") || mtype.Is("image/png")) {
		return nil, errImageType
	}

	return imageBytes, nil
}
