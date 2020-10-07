package models

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Event is an event
type Event struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Depth float64   `json:"depth"`
	Audit
}

// CreateEvent creates a new event
func CreateEvent(db *sqlx.DB, event *Event) (*Event, error) {
	var e Event
	if err := db.Get(
		&e, `INSERT INTO event (name, depth, creator, create_date) VALUES ($1, $2, $3, $4) RETURNING *`,
		event.Name, event.Depth, event.Creator, event.CreateDate,
	); err != nil {
		return nil, err
	}
	return &e, nil
}

// ListEvents lists all events
func ListEvents(db *sqlx.DB) ([]Event, error) {
	ee := make([]Event, 0)
	if err := db.Select(
		&ee, `SELECT * FROM event`,
	); err != nil {
		return make([]Event, 0), err
	}
	return ee, nil
}

// DeleteEvent deletes an event
func DeleteEvent(db *sqlx.DB, eventID *uuid.UUID) error {
	if _, err := db.Exec("DELETE FROM event WHERE id = $1", eventID); err != nil {
		return err
	}
	return nil
}
