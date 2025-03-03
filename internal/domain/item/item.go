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

type ItemAPI struct {
	ID           uint
	BrandId      uint
	BrandName    string
	Name         string
	Description  string
	Sex          int
	CategoryId   int
	CategoryType int
	CategoryName string
	Price        int
	Discount     *int
	OuterLink    string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	ImageId      []string
}

type ItemUpdate struct {
	ID          int
	BrandId     *int
	Name        *string
	Description *string
	Sex         *int
	CategoryId  *int
	Price       *uint
	Discount    *uint
	OuterLink   *string
}

type ItemCreate struct {
	BrandId     int
	Name        string
	Description string
	Sex         int
	CategoryId  int
	Price       uint
	Discount    uint
	OuterLink   string
}

type ItemInputData struct {
	ID         *uint
	BrandId    *uint
	Name       *string
	Sex        *int
	CategoryId *uint
	MinPrice   *uint
	MaxPrice   *uint
	Discount   *uint
	Offset     *uint
	Limit      *uint
}
