package domain

import "time"

// image model table image
type Image struct {
	ID         int
	ItemId     int
	ObjectId   string
	UploadedAt time.Time
}

type TempImage struct {
	ID         uint
	ObjectId   string
	UploadedAt time.Time
}
