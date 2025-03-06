package domain

import "time"

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
