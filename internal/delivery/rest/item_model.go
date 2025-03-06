package rest

type ItemQueryParams struct {
	ID         *uint   `query:"id"`
	BrandId    *uint   `query:"brand_id"`
	Name       *string `query:"name"`
	Sex        *int    `query:"sex"`
	CategoryId *uint   `query:"category_id"`
	MinPrice   *uint   `query:"min_price"`
	MaxPrice   *uint   `query:"max_price"`
	Discount   *uint   `query:"discount"`
	Offset     *uint   `query:"offset"`
	Limit      *uint   `query:"limit"`
}

type ItemUpdate struct {
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

type ItemCreate struct {
	BrandId     int    `json:"brand_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sex         int    `json:"sex" validate:"required"`
	CategoryId  int    `json:"category_id" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	Discount    uint   `json:"discount"`
	OuterLink   string `json:"outer_link" validate:"required"`
}
