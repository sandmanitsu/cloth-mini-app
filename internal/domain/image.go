package domain

type Image struct {
	ID     int    `json:"image_id"`
	ItemId int    `json:"item_id"`
	Link   string `json:"link"`
}
