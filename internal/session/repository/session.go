package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

	// GetSessionsByTable retrieves all sessions for a specific table
	GetSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error)

	// GetActiveSessionsByTable retrieves only active sessions for a specific table
	GetActiveSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error)

	// DeleteSession deletes a session by ID
	DeleteSession(ctx context.Context, id uuid.UUID) error

	// Table operations
	CreateTable(ctx context.Context, table *models.CreateTableRequest) (*models.Table, error)
	GetTable(ctx context.Context, id int) (*models.Table, error)
	GetTableByNumber(ctx context.Context, number int) (*models.Table, error)
	ListTables(ctx context.Context) ([]*models.Table, error)
	DeleteTable(ctx context.Context, id int) error
	BulkCreateTables(ctx context.Context, tableIDs []int) error
	TableExists(ctx context.Context, number int) (bool, error)
	IsTableAvailable(ctx context.Context, tableID int) (bool, error)
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
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET table_id = $1 WHERE id = $2", tableNumber, id)
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

// GetSessionsByTable retrieves all sessions for a specific table
func (r *postgresRepository) GetSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, table_id, created_at, completed_at, status FROM sessions WHERE table_id = $1 ORDER BY created_at DESC", tableID)
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

// GetActiveSessionsByTable retrieves only active sessions for a specific table
func (r *postgresRepository) GetActiveSessionsByTable(ctx context.Context, tableID int) ([]*models.Session, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, table_id, created_at, completed_at, status FROM sessions WHERE table_id = $1 AND status = $2 ORDER BY created_at DESC", tableID, models.StatusActive)
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

// DeleteSession deletes a session by ID
func (r *postgresRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}

// CreateTable creates a new table in the database
func (r *postgresRepository) CreateTable(ctx context.Context, req *models.CreateTableRequest) (*models.Table, error) {
	// Check if table number already exists
	exists, err := r.TableExists(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check table existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("table with id %d already exists", req.ID)
	}

	query := `INSERT INTO tables (id) VALUES ($1)`

	_, err = r.db.ExecContext(ctx, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &models.Table{ID: req.ID}, nil
}

// GetTable retrieves a table by ID
func (r *postgresRepository) GetTable(ctx context.Context, id int) (*models.Table, error) {
	query := `SELECT id FROM tables WHERE id = $1`

	var tableID int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&tableID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("table with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	return &models.Table{ID: tableID}, nil
}

// GetTableByNumber retrieves a table by table number
func (r *postgresRepository) GetTableByNumber(ctx context.Context, number int) (*models.Table, error) {
	return r.GetTable(ctx, number)
}

// ListTables lists all tables
func (r *postgresRepository) ListTables(ctx context.Context) ([]*models.Table, error) {
	query := `SELECT id FROM tables ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		var tableID int
		err := rows.Scan(&tableID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		tables = append(tables, &models.Table{ID: tableID})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	return tables, nil
}

// DeleteTable deletes a table by ID
func (r *postgresRepository) DeleteTable(ctx context.Context, id int) error {
	// Check if table has active sessions
	sessionQuery := `SELECT COUNT(*) FROM sessions WHERE table_id = $1 AND status = 'active'`
	var activeSessions int
	err := r.db.QueryRowContext(ctx, sessionQuery, id).Scan(&activeSessions)
	if err != nil {
		return fmt.Errorf("failed to check active sessions: %w", err)
	}
	if activeSessions > 0 {
		return fmt.Errorf("cannot delete table with active sessions")
	}

	query := `DELETE FROM tables WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("table with id %d not found", id)
	}

	return nil
}

// TableExists checks if a table with the given number exists
func (r *postgresRepository) TableExists(ctx context.Context, number int) (bool, error) {
	query := `SELECT COUNT(*) FROM tables WHERE id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, number).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}
	return count > 0, nil
}

// BulkCreateTables creates multiple tables in a single query
func (r *postgresRepository) BulkCreateTables(ctx context.Context, tableIDs []int) error {
	if len(tableIDs) == 0 {
		return nil
	}

	// Build the VALUES clause
	valueStrings := make([]string, 0, len(tableIDs))
	valueArgs := make([]interface{}, 0, len(tableIDs))
	for i, id := range tableIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i+1))
		valueArgs = append(valueArgs, id)
	}

	query := fmt.Sprintf("INSERT INTO tables (id) VALUES %s ON CONFLICT (id) DO NOTHING", strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	return err
}

// IsTableAvailable checks if a table has no active or pending sessions
func (r *postgresRepository) IsTableAvailable(ctx context.Context, tableID int) (bool, error) {
	query := `SELECT NOT EXISTS (
		SELECT 1 FROM sessions
		WHERE table_id = $1
			AND status IN ('active', 'pending')
	)`
	var available bool
	err := r.db.QueryRowContext(ctx, query, tableID).Scan(&available)
	if err != nil {
		return false, fmt.Errorf("failed to check table availability: %w", err)
	}
	return available, nil
}
