package repository

import (
	"context"
	"time"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type SessionRepository interface {
	Create(
		ctx context.Context,
		agentID string,
		name string,
		keyHash string,
		limit float64,
		expiresAt time.Time,
	) (*models.Session, error)

	FindByKeyHash(
		ctx context.Context,
		keyHash string,
	) (*models.Session, error)

	FindByAgentID(
		ctx context.Context,
		agentID string,
	) ([]models.Session, error)

	Revoke(
		ctx context.Context,
		id string,
	) error

	AddSpent(
		ctx context.Context,
		id string,
		amount float64,
	) error
}

type sessionRepository struct {
	db *pgx.Conn
}

func NewSessionRepository(
	db *pgx.Conn,
) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(
	ctx context.Context,
	agentID string,
	name string,
	keyHash string,
	limit float64,
	expiresAt time.Time,
) (*models.Session, error) {

	var session models.Session
	
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO agent_sessions (
			agent_id,
			name,
			key_hash,
			spending_limit,
			expires_at
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING
			id,
			agent_id,
			name,
			status,
			spending_limit,
			spent,
			created_at,
			expires_at
		`,
		agentID,
		name,
		keyHash,
		limit,
		expiresAt,
	).Scan(
		&session.ID,
		&session.AgentID,
		&session.Name,
		&session.Status,
		&session.Limit,
		&session.Spent,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) FindByKeyHash(
	ctx context.Context,
	keyHash string,
) (*models.Session, error) {

	var session models.Session

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			agent_id,
			name,
			status,
			spending_limit,
			spent,
			created_at,
			expires_at
		FROM agent_sessions
		WHERE key_hash = $1
		`,
		keyHash,
	).Scan(
		&session.ID,
		&session.AgentID,
		&session.Name,
		&session.Status,
		&session.Limit,
		&session.Spent,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) FindByAgentID(
	ctx context.Context,
	agentID string,
) ([]models.Session, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			agent_id,
			name,
			status,
			spending_limit,
			spent,
			created_at,
			expires_at
		FROM agent_sessions
		WHERE agent_id = $1
		ORDER BY created_at DESC
		`,
		agentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]models.Session, 0)

	for rows.Next() {

		var session models.Session

		err := rows.Scan(
			&session.ID,
			&session.AgentID,
			&session.Name,
			&session.Status,
			&session.Limit,
			&session.Spent,
			&session.CreatedAt,
			&session.ExpiresAt,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func (r *sessionRepository) Revoke(
	ctx context.Context,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE agent_sessions
		SET
			status = 'revoked',
			revoked_at = NOW()
		WHERE id = $1
		AND status = 'active'
		`,
		id,
	)

	return err
}

func (r *sessionRepository) AddSpent(
	ctx context.Context,
	id string,
	amount float64,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE agent_sessions
		SET spent = spent + $1
		WHERE id = $2
		`,
		amount,
		id,
	)

	return err
}
