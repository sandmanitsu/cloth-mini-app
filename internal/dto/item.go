package dto

type ItemDTO struct {
	ID          int     `param:"id"`
	BrandId     *int    `json:"brand_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Sex         *int    `json:"sex"`
	CategoryId  *int    `json:"category_id"`
	Price       *uint   `json:"price"`
	Discount    *uint   `json:"discount"`
	OuterLink   *string `json:"outerlink"`
}

type ItemCreateDTO struct {
	BrandId     int    `json:"brand_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sex         int    `json:"sex" validate:"required"`
	CategoryId  int    `json:"category_id" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	Discount    uint   `json:"discount"`
	OuterLink   string `json:"outer_link" validate:"required"`
}
