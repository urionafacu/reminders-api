package main

import (
	"fmt"
	"log"
	"net/http"
	"reminders-api/internal/api"
	"reminders-api/internal/config"
	"reminders-api/internal/service"
	"reminders-api/internal/storage"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := storage.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize the database
	err = db.Init()
	if err != nil {
		log.Fatalf("Error initializing database schema: %v", err)
	}

	eventService := service.NewEventService(db)
	emailService := service.NewEmailService(cfg)

	router := api.SetupRoutes(eventService, emailService, cfg.APIKey)

	// Start the reminder service in a separate goroutine
	go runDailyReminders(eventService, emailService)

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

func runDailyReminders(eventService *service.EventService, emailService *service.EmailService) {
	for {
		events, err := eventService.GetAllEvents()
		if err != nil {
			log.Printf("Error getting events for reminders: %v", err)
		} else {
			emailService.SendDailyReminders(events)
		}

		// Wait until the next execution (every 24 hours)
		time.Sleep(24 * time.Hour)
	}
}
