package models

import (
	"encoding/json"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-simple-asyncer/asyncer"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Compute struct {
	ID         uuid.UUID `json:"id"`
	EventID    uuid.UUID `json:"event_id" db:"event_id"`
	EventDepth float64   `json:"event_depth" db:"event_depth"`
	FIPS       string    `json:"fips"`
	Audit
}

func GetCompute(db *sqlx.DB, id *uuid.UUID) (*Compute, error) {
	var c Compute
	if err := db.Get(
		&c, `SELECT * FROM v_compute WHERE id=$1`, id,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func RunCompute(db *sqlx.DB, ae asyncer.Asyncer, c *Compute) (*Compute, error) {
	var id uuid.UUID
	if err := db.Get(
		&id,
		`INSERT INTO compute (event_id, fips, creator, create_date) VALUES ($1,$2,$3,$4)
		 RETURNING ID`, c.EventID, c.FIPS, c.Creator, c.CreateDate,
	); err != nil {
		return nil, err
	}
	nc, err := GetCompute(db, &id)
	if err != nil {
		return nil, err
	}
	// Build go-consequences FipsCodeCompute; Convert to ByteString
	payload, err := json.Marshal(
		compute.RequestArgs{
			Args: compute.FipsCodeCompute{
				ID:         nc.ID.String(),
				FIPS:       nc.FIPS,
				HazardArgs: hazards.DepthEvent{Depth: nc.EventDepth},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	// Send the Message
	if err := ae.CallAsync("", payload); err != nil {
		return nil, err
	}
	return nc, nil
}
