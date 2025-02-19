package domain

type Category struct {
	ID   int    `json:"id"`
	Type int    `json:"type"`
	Name string `json:"category_name"`
}
