//go:build integration

package integrations

import (
	"bytes"
	domain "cloth-mini-app/internal/domain/item"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type GetItem struct {
	ID           uint       `json:"id"`
	BrandId      uint       `json:"brand_id"`
	BrandName    string     `json:"brand_name"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Sex          int        `json:"sex"`
	CategoryId   int        `json:"category_id"`
	CategoryType int        `json:"category_type"`
	CategoryName string     `json:"category_name"`
	Price        int        `json:"price"`
	Discount     *int       `json:"discount"`
	OuterLink    string     `json:"outer_link"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type ItemResponse struct {
	Count int       `json:"count"`
	Items []GetItem `json:"items"`
}

func (i *IntegrationSuite) TestGetItem() {
	url := host + "/item/get?limit=10&offset=0"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	i.Require().Equal(http.StatusOK, response.StatusCode)

	respItem, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var items ItemResponse
	err = json.Unmarshal(respItem, &items)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(3, items.Count)
}

type ItemByIdResponse struct {
	ID           uint       `json:"id"`
	BrandId      uint       `json:"brand_id"`
	BrandName    string     `json:"brand_name"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Sex          int        `json:"sex"`
	CategoryId   int        `json:"category_id"`
	CategoryType int        `json:"category_type"`
	CategoryName string     `json:"category_name"`
	Price        int        `json:"price"`
	Discount     *int       `json:"discount"`
	OuterLink    string     `json:"outer_link"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	ImageId      []string   `json:"image_id"`
}

func (i *IntegrationSuite) TestGetItemById() {
	item := domain.ItemCreate{
		BrandId:     1,
		Name:        "test item by id",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	id := i.createItem(item)
	itemId := strconv.Itoa(int(id))
	url := host + "/item/get/" + itemId

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	i.Require().Equal(http.StatusOK, response.StatusCode)

	rawItem, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var respItem ItemByIdResponse
	err = json.Unmarshal(rawItem, &respItem)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(item.Name, respItem.Name)
	i.Require().Equal(item.Price, uint(respItem.Price))
}

func (i *IntegrationSuite) createItem(item domain.ItemCreate) uint {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("items").
		Columns("brand_id", "name", "description", "sex", "category_id", "price", "discount", "outer_link", "created_at").
		Values(item.BrandId, item.Name, item.Description, item.Sex, item.CategoryId, item.Price, item.Discount, item.OuterLink, time.Now()).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	var itemId uint
	err = i.db.QueryRow(sql, args...).Scan(&itemId)
	if err != nil {
		log.Fatal(err)
	}

	return itemId
}

type ItemUpdateField struct {
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

func (i *IntegrationSuite) TestUpdateItem() {
	item := domain.ItemCreate{
		BrandId:     1,
		Name:        "test update item",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	id := i.createItem(item)
	itemId := strconv.Itoa(int(id))
	url := host + "/item/update/" + itemId

	updateField := ItemUpdateField{
		Name:  "updated",
		Price: 1234,
	}

	json, err := json.Marshal(updateField)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	i.Require().Equal(http.StatusOK, response.StatusCode)

	dbItem, err := i.getItem(id)
	i.Require().NoError(err)
	i.Require().Equal(updateField.Name, dbItem.Name)
	i.Require().Equal(updateField.Price, dbItem.Price)
}

func (i *IntegrationSuite) getItem(itemId uint) (domain.ItemCreate, error) {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"i.id", "i.brand_id", "i.name", "i.description", "i.sex", "i.category_id",
			"i.price", "i.discount", "i.outer_link",
		).
		From("items i").
		Where("id = ?", itemId).
		ToSql()
	if err != nil {
		return domain.ItemCreate{}, err
	}

	var item domain.ItemCreate
	err = i.db.QueryRow(sql, args...).Scan(
		&itemId,
		&item.BrandId,
		&item.Name,
		&item.Description,
		&item.Sex,
		&item.CategoryId,
		&item.Price,
		&item.Discount,
		&item.OuterLink,
	)
	if err != nil {
		return domain.ItemCreate{}, err
	}

	return item, nil
}

func (i *IntegrationSuite) TestCreateItem() {
	url := host + "/item/create"

	item := ItemCreate{
		BrandId:     1,
		Name:        "create test",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images: []string{
			uuid.NewString(),
			uuid.NewString(),
		},
	}

	json, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	i.Require().Equal(http.StatusOK, response.StatusCode)

	dbItem := i.getCreateItem(item.Name)

	i.Require().Equal(item.Name, dbItem.Name)
	i.Require().ElementsMatch(item.Images, dbItem.Images)
	i.Require().Equal(item.BrandId, dbItem.BrandId)
	i.Require().Equal(item.CategoryId, dbItem.CategoryId)
}

func (i *IntegrationSuite) getCreateItem(name string) domain.ItemCreate {
	const op = "getCreateItem"

	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"i.id", "i.brand_id", "i.name", "i.description", "i.sex", "i.category_id",
			"i.price", "i.discount", "i.outer_link",
		).
		From("items i").
		Where("i.name = ?", name).
		ToSql()
	if err != nil {
		log.Fatal(op, err)
	}

	var itemId int
	var item domain.ItemCreate
	err = i.db.QueryRow(sql, args...).Scan(
		&itemId,
		&item.BrandId,
		&item.Name,
		&item.Description,
		&item.Sex,
		&item.CategoryId,
		&item.Price,
		&item.Discount,
		&item.OuterLink,
	)
	if err != nil {
		log.Fatal(op, err)
	}

	item.Images = i.getImages(itemId)

	return item
}

func (i *IntegrationSuite) getImages(itemId int) []string {
	sql, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("object_id").
		From("images").
		Where("item_id = ?", itemId).
		ToSql()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := i.db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			log.Fatal(err)
		}

		images = append(images, image)
	}

	return images
}

func (i *IntegrationSuite) TestDeleteItem() {
	item := domain.ItemCreate{
		BrandId:     1,
		Name:        "test delete",
		Description: "some description...",
		Sex:         1,
		CategoryId:  1,
		Price:       10000,
		Discount:    10,
		OuterLink:   "http:/localhost:8080/",
		Images:      nil,
	}

	id := i.createItem(item)
	itemId := strconv.Itoa(int(id))
	url := host + "/item/delete/" + itemId

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

	_, err = i.getItem(id)
	i.Require().Error(err)
}
