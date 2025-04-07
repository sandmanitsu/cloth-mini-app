//go:build integration

package integrations

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Category struct {
	CategoryId int    `json:"category_id"`
	Type       int    `json:"type"`
	Name       string `json:"category_name"`
}

func (i *IntegrationSuite) TestGetCategories() {
	url := host + "/category/get"

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

	respCat, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var categs []Category
	err = json.Unmarshal(respCat, &categs)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Greater(len(categs), 0)
}
