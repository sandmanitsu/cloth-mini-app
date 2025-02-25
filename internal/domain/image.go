package domain

// image model table image
type Image struct {
	ID       int    `json:"image_id"`
	ItemId   int    `json:"item_id"`
	ObjectId string `json:"object_id"`
}
