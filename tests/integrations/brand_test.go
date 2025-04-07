//go:build integration

package integrations

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Brand struct {
	ID   int    `json:"brand_id"`
	Name string `json:"brand_name"`
}

func (i *IntegrationSuite) TestGetBrands() {
	url := host + "/brand/get"

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

	respBrands, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var brands []Brand
	err = json.Unmarshal(respBrands, &brands)
	if err != nil {
		log.Fatal(err)
	}

	i.Require().Equal(3, len(brands))
}
