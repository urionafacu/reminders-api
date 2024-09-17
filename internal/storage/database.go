package storage

import (
	"reminders-api/internal/models"
)

// Database define el contrato para las operaciones de base de datos
type Database interface {
	// Events
	GetEvents() ([]models.Event, error)
	GetEventByID(id int64) (*models.Event, error)
	AddEvent(event *models.Event) error
	UpdateEvent(event *models.Event) error
	DeleteEvent(id int64) error

	Init() error
	Close() error
}
