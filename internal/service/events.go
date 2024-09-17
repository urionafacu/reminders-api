package service

import (
	"reminders-api/internal/models"
	"reminders-api/internal/storage"
)

type EventService struct {
	db storage.Database
}

func NewEventService(db storage.Database) *EventService {
	return &EventService{db: db}
}

func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.db.GetEvents()
}

func (s *EventService) AddEvent(event *models.Event) error {
	return s.db.AddEvent(event)
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	return s.db.UpdateEvent(event)
}

func (s *EventService) DeleteEvent(eventID int64) error {
	return s.db.DeleteEvent(eventID)
}

func (s *EventService) GetEventByID(eventID int64) (*models.Event, error) {
	return s.db.GetEventByID(eventID)
}
