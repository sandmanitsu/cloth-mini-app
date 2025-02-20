package domain

import "time"

type Item struct {
	ID          uint       `json:"id" gorm:"primarykey"`
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
	ID           uint       `json:"id" gorm:"primarykey"`
	Brand        string     `json:"brand"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Sex          int        `json:"sex"`
	CategoryType int        `json:"category_type"`
	CategoryName string     `json:"category_name"`
	Price        int        `json:"price"`
	Discount     int        `json:"discount"`
	OuterLink    string     `json:"outer_link"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
