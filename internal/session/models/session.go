package models

import (
	"time"

	"github.com/google/uuid"
)

// SessionStatus represents the possible states of a session
type SessionStatus string

const (
	StatusActive    SessionStatus = "active"
	StatusCompleted SessionStatus = "completed"
	StatusPending   SessionStatus = "pending"
	StatusCancelled SessionStatus = "cancelled"
)

// Session represents a dining session started by QR scan
type Session struct {
	ID          uuid.UUID     // unique session ID (UUID)
	TableID     int           // which table this session is for
	CreatedAt   time.Time     // when the session was created
	CompletedAt *time.Time    // when the session was completed, nil if not completed
	Status      SessionStatus // e.g., StatusActive, StatusCompleted, or StatusPending
}

type Bill struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	Total     float64
	Subtotal  float64
	Tax       float64
	CreatedAt time.Time
}

// Table represents a physical table in the restaurant
type Table struct {
	ID int `json:"id" db:"id"` // table number (primary key)
}

// CreateTableRequest represents the request to create a new table
type CreateTableRequest struct {
	ID int `json:"id" validate:"required,min=1" db:"id"`
}

// UpdateTableRequest represents the request to update an existing table
type UpdateTableRequest struct {
	ID int `json:"id" validate:"required,min=1" db:"id"`
}
