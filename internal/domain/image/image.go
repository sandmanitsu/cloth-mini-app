package domain

import "time"

// image model table image
type Image struct {
	ID         int       `json:"image_id"`
	ItemId     int       `json:"item_id"`
	ObjectId   string    `json:"object_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}
