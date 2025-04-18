package domain

import "time"

const (
	EventCreateItem = "create_item"
)

type Event struct {
	Id         int
	EventType  string
	Payload    []byte
	Status     string
	CreatedAt  time.Time
	ReservedTo *time.Time
}
