package models

import "time"

// Audit captures information on create and modify
type Audit struct {
	Creator    int       `json:"creator"`
	CreateDate time.Time `json:"create_date" db:"create_date"`
}