package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"restaurant/internal/session/models"

	"github.com/google/uuid"
)

// Repository defines methods for session database operations
type Repository interface {
	// CreateSession creates a new session with the given ID and table ID
	CreateSession(ctx context.Context, id uuid.UUID, tableID int) (*models.Session, error)

	// GetSession retrieves a session by ID
	GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error)

	// UpdateSession updates the status of a session
	UpdateSession(ctx context.Context, id uuid.UUID, newStatus models.SessionStatus) error

	// ListSessions lists sessions with pagination (offset and limit)
	ListSessions(ctx context.Context, offset int, limit int) ([]*models.Session, error)

	// ListActiveSessions lists all sessions with status "active"
	ListActiveSessions(ctx context.Context) ([]*models.Session, error)

	// ChangeSessionTable changes the table ID of a session by table number
	ChangeSessionTable(ctx context.Context, id uuid.UUID, tableNumber int) error
}

// postgresRepository implements Repository using PostgreSQL
type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL-based repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// UpdateSession updates the status of a session in the database
func (r *postgresRepository) UpdateSession(
	ctx context.Context,
	id uuid.UUID,
	newStatus models.SessionStatus,
) error {
	if newStatus == models.StatusCompleted {
		now := time.Now()
		_, err := r.db.ExecContext(ctx,
			"UPDATE sessions SET status = $1, completed_at = $2 WHERE id = $3",
			newStatus, now, id,
		)
		return err
	}

	_, err := r.db.ExecContext(ctx,
		"UPDATE sessions SET status = $1, completed_at = NULL WHERE id = $2",
		newStatus, id,
	)
	return err
}

// ListSessions retrieves a paginated list of sessions from the database
func (r *postgresRepository) ListSessions(ctx context.Context, offset int, limit int) ([]*models.Session, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, table_id, created_at, completed_at, status FROM sessions ORDER BY created_at DESC OFFSET $1 LIMIT $2", offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		var status string
		var completedAt *time.Time
		err := rows.Scan(&session.ID, &session.TableID, &session.CreatedAt, &completedAt, &status)
		if err != nil {
			return nil, err
		}
		session.Status = models.SessionStatus(status)
		session.CompletedAt = completedAt
		sessions = append(sessions, &session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

// ListActiveSessions retrieves all sessions with status "active"
func (r *postgresRepository) ListActiveSessions(ctx context.Context) ([]*models.Session, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, table_id, created_at, completed_at, status FROM sessions WHERE status = $1", models.StatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		var status string
		var completedAt *time.Time
		err := rows.Scan(&session.ID, &session.TableID, &session.CreatedAt, &completedAt, &status)
		if err != nil {
			return nil, err
		}
		session.Status = models.SessionStatus(status)
		session.CompletedAt = completedAt
		sessions = append(sessions, &session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

// ChangeSessionTable changes the table ID of a session by table number
func (r *postgresRepository) ChangeSessionTable(ctx context.Context, id uuid.UUID, tableNumber int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET table_id = (SELECT id FROM tables WHERE number = $1) WHERE id = $2", tableNumber, id)
	return err
}

// CreateSession inserts a new session into the database
func (r *postgresRepository) CreateSession(ctx context.Context, id uuid.UUID, tableID int) (*models.Session, error) {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, "INSERT INTO sessions (id, table_id, created_at, completed_at, status) VALUES ($1, $2, $3, $4, $5)", id, tableID, now, nil, models.StatusActive)
	if err != nil {
		return nil, err
	}
	// Return the created session
	return &models.Session{
		ID:          id,
		TableID:     tableID,
		CreatedAt:   now,
		CompletedAt: nil,
		Status:      models.StatusActive,
	}, nil
}

// GetSession retrieves a session by ID from the database
func (r *postgresRepository) GetSession(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	var status string
	var completedAt *time.Time
	err := r.db.QueryRowContext(ctx, "SELECT id, table_id, created_at, completed_at, status FROM sessions WHERE id = $1", id).Scan(
		&session.ID, &session.TableID, &session.CreatedAt, &completedAt, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Session not found") // or return an error like errors.New("session not found")
		}
		return nil, err
	}
	session.Status = models.SessionStatus(status)
	session.CompletedAt = completedAt
	return &session, nil
}
