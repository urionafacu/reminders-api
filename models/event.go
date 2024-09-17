package models

import (
	"time"
)

type Event struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Recurring bool      `json:"recurring"`
}
