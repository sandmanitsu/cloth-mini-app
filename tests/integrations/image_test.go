package integrations

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type ItemCreate struct {
	BrandId     int      `json:"brand_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Sex         int      `json:"sex"`
	CategoryId  int      `json:"category_id"`
	Price       uint     `json:"price"`
	Discount    uint     `json:"discount"`
	OuterLink   string   `json:"outer_link"`
	Images      []string `json:"temp_images"`
}

func (i *IntegrationSuite) TestCreateTempImage() {
	url := host + "/image/temp"

	image, err := os.Open("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer image.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", "test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, image)
	if err != nil {
		log.Fatal(err)
	}

	imageID := uuid.NewString()
	err = writer.WriteField("uuid", imageID)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(http.StatusOK, response.StatusCode)

	dbImageID := i.getTempImage(imageID)
	i.Require().Equal(imageID, dbImageID)

	minioFileId := i.getMinioFileId(imageID)
	i.Require().Equal(imageID, minioFileId)
}

func (i *IntegrationSuite) getTempImage(id string) string {
	const op = "image_test.getTempImage"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("object_id").
		From("temp_images").
		Where("object_id = ?", id).
		ToSql()
	if err != nil {
		log.Fatal(op, err)
	}

	var imageId string
	err = i.db.QueryRow(sql, args...).Scan(&imageId)
	if err != nil {
		log.Fatal(op, err)
	}

	return imageId
}

func (i *IntegrationSuite) TestCreateImage() {
	url := host + "/image/create?itemId=1"

	image, err := os.Open("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer image.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", "test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, image)
	if err != nil {
		log.Fatal(err)
	}

	imageID := uuid.NewString()
	err = writer.WriteField("uuid", imageID)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.StatusCode)
	i.Require().Equal(http.StatusOK, response.StatusCode)

	dbImageID, err := i.getImage(mockItemID)
	i.Require().NoError(err)

	minioFileId := i.getMinioFileId(dbImageID)
	i.Require().Equal(dbImageID, minioFileId)
}

func (i *IntegrationSuite) getImage(itemId int) (string, error) {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("object_id").
		From("images").
		Where("item_id = ?", itemId).
		ToSql()
	if err != nil {
		return "", err
	}

	var imageId string
	err = i.db.QueryRow(sql, args...).Scan(&imageId)
	if err != nil {
		return "", err
	}

	return imageId, nil
}

func (i *IntegrationSuite) TestGetImage() {
	imageId := uuid.NewString()
	url := host + "/image/get/" + imageId

	image, err := os.ReadFile("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}

	i.putImageToMinio(imageId, image)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(http.StatusOK, response.StatusCode)
	i.Require().Equal("image/jpeg", response.Header.Get("Content-Type"))

	responseImage, err := io.ReadAll(response.Body)
	if err != nil {
		i.Require().NoError(err)
	}
	i.Require().True(reflect.DeepEqual(image, responseImage))
}

func (i *IntegrationSuite) TestDeleteImage() {
	imageId := uuid.NewString()
	url := host + "/image/delete?image_id=" + imageId

	image, err := os.ReadFile("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}

	i.putImageToMinio(imageId, image)
	i.createImageDB(1, imageId)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(http.StatusOK, response.StatusCode)

	_, err = i.getImage(mockItemID)
	i.Require().Error(err)
}

func (i *IntegrationSuite) createImageDB(itemId int, objectId string) {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("images").
		Columns("item_id", "object_id", "uploaded_at").
		Values(itemId, objectId, time.Now()).
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	_, err = i.db.Exec(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *IntegrationSuite) TestMaxImagePerItem() {
	url := host + "/image/create?itemId=1"

	maxImagesPerItem := 4
	for range maxImagesPerItem {
		i.createImageDB(mockItemID, uuid.NewString())
	}

	image, err := os.Open("fixtures/test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer image.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", "test_pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, image)
	if err != nil {
		log.Fatal(err)
	}

	imageID := uuid.NewString()
	err = writer.WriteField("uuid", imageID)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(http.StatusBadRequest, response.StatusCode)
}
