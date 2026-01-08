package session

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
	ID          uuid.UUID     `json:"id"`           // unique session ID (UUID)
	TableID     int           `json:"table_id"`     // which table this session is for
	CreatedAt   time.Time     `json:"created_at"`   // when the session was created
	CompletedAt *time.Time    `json:"completed_at"` // when the session was completed, nil if not completed
	Status      SessionStatus `json:"status"`       // e.g., StatusActive, StatusCompleted, or StatusPending
}

type Bill struct {
	ID        uuid.UUID `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	Total     float64   `json:"total"`
	Subtotal  float64   `json:"subtotal"`
	Tax       float64   `json:"tax"`
	CreatedAt time.Time `json:"created_at"`
}

// Table represents a physical table in the restaurant
type Table struct {
	ID int `json:"id" db:"id"` // table number (primary key)
}
