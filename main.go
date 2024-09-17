package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reminders-api/config"
	"reminders-api/db"
	"reminders-api/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
)

var (
	database *db.DB
	cfg      *config.Config
)

func sendEmail(config *config.Config, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.EmailFrom)
	m.SetHeader("To", config.EmailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.EmailFrom, config.EmailPassword)

	return d.DialAndSend(m)
}

func authenticateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != cfg.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := database.GetEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(events)
}

func addEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := database.AddEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event models.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	event.ID = id

	err = database.UpdateEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	cfg = config.LoadConfig()

	if cfg.EmailFrom == "" || cfg.EmailTo == "" || cfg.EmailPassword == "" || cfg.APIKey == "" {
		log.Fatal("Missing essential configuration. Please check your .env file or environment variables.")
	}

	var err error
	database, err = db.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	err = database.CreateTable()
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	events, err := database.GetEvents()
	if err != nil {
		log.Fatalf("Error getting events: %v", err)
	}

	if len(events) == 0 {
		err = database.LoadInitialEvents("initial_events.json")
		if err != nil {
			log.Printf("Error loading initial events: %v", err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/events", authenticateMiddleware(getEventsHandler)).Methods("GET").Headers("Content-Type", "application/json")
	r.HandleFunc("/events", authenticateMiddleware(addEventHandler)).Methods("POST")
	r.HandleFunc("/events/{id}", authenticateMiddleware(updateEventHandler)).Methods("PUT")
	r.HandleFunc("/events/{id}", authenticateMiddleware(deleteEventHandler)).Methods("DELETE")

	go func() {
		for {
			events, err := database.GetEvents()
			if err != nil {
				log.Printf("Error getting events: %v", err)
				time.Sleep(24 * time.Hour)
				continue
			}

			now := time.Now()
			for _, event := range events {
				if event.Recurring {
					if now.Day() == event.Date.Day() {
						subject := fmt.Sprintf("Reminder: %s", event.Name)
						body := fmt.Sprintf("Today is your %s!", event.Name)
						err := sendEmail(cfg, subject, body)
						if err != nil {
							log.Printf("Error sending email for %s: %v\n", event.Name, err)
						}
					}
				} else {
					daysUntil := event.Date.Sub(now).Hours() / 24
					if daysUntil <= 7 && daysUntil > 0 {
						subject := fmt.Sprintf("Reminder: %s in %d days", event.Name, int(daysUntil))
						body := fmt.Sprintf("%d days until %s!", int(daysUntil), event.Name)
						err := sendEmail(cfg, subject, body)
						if err != nil {
							log.Printf("Error sending email for %s: %v\n", event.Name, err)
						}
					}
				}
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
