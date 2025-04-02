package integrations

import (
	domain "cloth-mini-app/internal/domain/item"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
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
	url := "http://localhost:8080/item/get?limit=10&offset=0"

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
	url := "http://localhost:8080/item/get/" + itemId

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
