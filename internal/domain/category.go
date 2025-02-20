package domain

// Category model table category
type Category struct {
	CategoryId int    `json:"category_id"`
	Type       int    `json:"type"`
	Name       string `json:"category_name"`
}
