package api

import (
	"reminders-api/internal/api/handlers"
	"reminders-api/internal/api/middleware"
	"reminders-api/internal/service"

	"github.com/gorilla/mux"
)

func SetupRoutes(eventService *service.EventService, emailService *service.EmailService, apiKey string) *mux.Router {
	r := mux.NewRouter()

	eventHandler := handlers.NewEventHandler(eventService)

	r.Use(middleware.APIKeyAuth(apiKey))

	r.HandleFunc("/events", eventHandler.GetEvents).Methods("GET")
	r.HandleFunc("/events", eventHandler.AddEvent).Methods("POST")
	r.HandleFunc("/events/{id}", eventHandler.UpdateEvent).Methods("PUT")
	r.HandleFunc("/events/{id}", eventHandler.DeleteEvent).Methods("DELETE")

	return r
}
