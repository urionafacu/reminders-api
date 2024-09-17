package storage

import (
	"database/sql"
	"reminders-api/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dataSourceName string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TEXT NOT NULL,
		recurring BOOLEAN NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteDB) GetEvents() ([]models.Event, error) {
	rows, err := s.db.Query("SELECT id, name, date, recurring FROM events")
	if err != nil {
		return []models.Event{}, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		var e models.Event
		var dateStr string
		err := rows.Scan(&e.ID, &e.Name, &dateStr, &e.Recurring)
		if err != nil {
			return []models.Event{}, err
		}
		e.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return []models.Event{}, err
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return []models.Event{}, err
	}

	return events, nil
}

func (s *SQLiteDB) GetEventByID(id int64) (*models.Event, error) {
	var e models.Event
	var dateStr string
	err := s.db.QueryRow("SELECT id, name, date, recurring FROM events WHERE id = ?", id).
		Scan(&e.ID, &e.Name, &dateStr, &e.Recurring)
	if err != nil {
		return nil, err
	}
	e.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *SQLiteDB) AddEvent(event *models.Event) error {
	result, err := s.db.Exec("INSERT INTO events(name, date, recurring) VALUES(?, ?, ?)",
		event.Name, event.Date.Format(time.RFC3339), event.Recurring)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = id
	return nil
}

func (s *SQLiteDB) UpdateEvent(event *models.Event) error {
	_, err := s.db.Exec("UPDATE events SET name = ?, date = ?, recurring = ? WHERE id = ?",
		event.Name, event.Date.Format(time.RFC3339), event.Recurring, event.ID)
	return err
}

func (s *SQLiteDB) DeleteEvent(id int64) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = ?", id)
	return err
}
