package validation

import (
	"errors"
	"restaurant/internal/session/models"

	"github.com/google/uuid"
)

// CreateSessionRequest represents the request to create a session
type CreateSessionRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

// UpdateSessionRequest represents the request to update a session
type UpdateSessionRequest struct {
	Status models.SessionStatus `json:"status" validate:"required,oneof=active completed pending cancelled"`
}

// ListSessionsRequest represents the request to list sessions with pagination
type ListSessionsRequest struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"required,min=1,max=100"`
}

// ChangeSessionTableRequest represents the request to change a session's table
type ChangeSessionTableRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

// ValidateCreateSession validates the create session request
func ValidateCreateSession(req CreateSessionRequest) error {
	return ValidateStruct(req)
}

// ValidateUpdateSession validates the update session request
func ValidateUpdateSession(req UpdateSessionRequest) error {
	return ValidateStruct(req)
}

// ValidateListSessions validates the list sessions request
func ValidateListSessions(req ListSessionsRequest) error {
	return ValidateStruct(req)
}

// ValidateChangeSessionTable validates the change session table request
func ValidateChangeSessionTable(req ChangeSessionTableRequest) error {
	return ValidateStruct(req)
}

// ValidateSessionID validates a session ID
func ValidateSessionID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid session ID")
	}
	return nil
}
