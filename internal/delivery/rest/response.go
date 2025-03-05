package rest

import "time"

type ErrorResponse struct {
	Err string `json:"error"`
}

type SuccessResponse struct {
	Status    bool   `json:"status"`
	Operation string `json:"operation"`
}

type ItemResponse struct {
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

type ItemsResponse struct {
	Count int            `json:"count"`
	Items []ItemResponse `json:"items"`
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
