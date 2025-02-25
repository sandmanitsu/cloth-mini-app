package domain

import "time"

// item model table items
type Item struct {
	ID          uint       `json:"id"`
	Brand       string     `json:"brand"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Sex         int        `json:"sex"`
	CategoryId  int        `json:"category_id"`
	Price       int        `json:"price"`
	Discount    *int       `json:"discount"`
	OuterLink   string     `json:"outer_link"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// item model with field from category table
type ItemAPI struct {
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
