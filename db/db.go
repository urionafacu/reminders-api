package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"reminders-api/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS events(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date DATETIME NOT NULL,
		recurring BOOLEAN NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func (db *DB) UpdateEvent(event models.Event) error {
	query := `UPDATE events SET name = ?, date = ?, recurring = ? WHERE id = ?`
	_, err := db.Exec(query, event.Name, event.Date, event.Recurring, event.ID)
	return err
}

func (db *DB) DeleteEvent(id int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func (db *DB) LoadInitialEvents(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var data struct {
		Events []struct {
			Name      string `json:"name"`
			Date      string `json:"date"`
			Recurring bool   `json:"recurring"`
		} `json:"events"`
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	for _, eventData := range data.Events {
		date, err := time.Parse(time.RFC3339, eventData.Date)
		if err != nil {
			log.Printf("Error parsing date for event %s: %v", eventData.Name, err)
			continue
		}

		event := models.Event{
			Name:      eventData.Name,
			Date:      date,
			Recurring: eventData.Recurring,
		}

		_, err = db.AddEvent(event)
		if err != nil {
			log.Printf("Error adding event %s: %v", event.Name, err)
		} else {
			log.Printf("Successfully added event: %s, Date: %v", event.Name, event.Date)
		}
	}

	return nil
}

func (db *DB) AddEvent(event models.Event) (int64, error) {
	query := `INSERT INTO events(name, date, recurring) VALUES(?, ?, ?)`
	result, err := db.Exec(query, event.Name, event.Date.Format(time.RFC3339), event.Recurring)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *DB) GetEvents() ([]models.Event, error) {
	rows, err := db.Query("SELECT id, name, date, recurring FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		var dateStr string
		err := rows.Scan(&e.ID, &e.Name, &dateStr, &e.Recurring)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		e.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			continue
		}
		events = append(events, e)
	}

	return events, nil
}
