package domain

type Item struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Brand       string `json:"brand"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sex         int    `json:"sex"`
	CategoryId  int    `json:"category_id"`
	SizeAmount  string `json:"size_amount"`
	Price       int    `json:"price"`
	Discount    int    `json:"discount"`
	OuterLink   string `json:"outer_link"`
}

type ItemAPI struct {
	ID           uint   `json:"id" gorm:"primarykey"`
	Brand        string `json:"brand"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Sex          int    `json:"sex"`
	CategoryType int    `json:"category_type"`
	CategoryName string `json:"category_name"`
	SizeAmount   string `json:"size_amount"`
	Price        int    `json:"price"`
	Discount     int    `json:"discount"`
	OuterLink    string `json:"outer_link"`
}
