package models

import (
	"time"
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
	ID          string       // unique session ID (UUID)
	TableID     int          // which table this session is for
	CreatedAt   time.Time    // when the session was created
	CompletedAt *time.Time   // when the session was completed, nil if not completed
	Status      SessionStatus // e.g., StatusActive, StatusCompleted, or StatusPending
}

type Table struct {
	ID        string    // unique table ID
	Number    int       // table number in the restaurant
	CreatedAt time.Time // when the table was created
}


type Bill struct {
    ID         string
    SessionID  string
    Total      float64
    Subtotal   float64
    Tax        float64
    CreatedAt  time.Time
}
