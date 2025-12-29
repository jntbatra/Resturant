package service

import (
	"restaurant/internal/session/models"
	"restaurant/internal/session/repository"
)

// Service defines business logic for sessions
type Service interface {
	CreateSession(id string, tableID int) (*models.Session, error)
	GetSession(id string) (*models.Session, error)
	UpdateSession(id string, status models.SessionStatus) (*models.Session, error)
	ListSessions(offset, limit int) ([]*models.Session, error)
	ListActiveSessions() ([]*models.Session, error)
	ChangeTable(id string, tableNumber int) error
}

// sessionService implements Service
type sessionService struct {
	repo repository.Repository
}

// NewService creates a new session service
func NewService(repo repository.Repository) Service {
	return &sessionService{repo: repo}
}

// CreateSession creates a new session
func (s *sessionService) CreateSession(id string, tableID int) (*models.Session, error) {
	return s.repo.CreateSession(id, tableID)
}

// GetSession retrieves a session
func (s *sessionService) GetSession(id string) (*models.Session, error) {
	return s.repo.GetSession(id)
}

// UpdateSession updates the status of a session
func (s *sessionService) UpdateSession(id string, status models.SessionStatus) (*models.Session, error) {
	err := s.repo.UpdateSession(id, status)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSession(id)
}

// ListSessions lists sessions with pagination
func (s *sessionService) ListSessions(offset, limit int) ([]*models.Session, error) {
	return s.repo.ListSessions(offset, limit)
}

// ListActiveSessions lists active sessions
func (s *sessionService) ListActiveSessions() ([]*models.Session, error) {
	return s.repo.ListActiveSessions()
}

// ChangeTable changes the table of a session
func (s *sessionService) ChangeTable(id string, tableNumber int) error {
	return s.repo.ChangeSessionTable(id, tableNumber)
}