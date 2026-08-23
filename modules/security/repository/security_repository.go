// modules/security/repository/security_repository.go

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type SecurityRepository interface {

	Create(
		ctx context.Context,
		agentID *string,
		eventType string,
		title string,
		description string,
		result string,
	) error
}

type securityRepository struct {
	db *pgx.Conn
}

func NewSecurityRepository(
	db *pgx.Conn,
) SecurityRepository {
	return &securityRepository{db: db}
}

func (r *securityRepository) Create(
	ctx context.Context,
	agentID *string,
	eventType string,
	title string,
	description string,
	result string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO security_events (
			agent_id,
			type,
			title,
			description,
			result
		)
		VALUES ($1,$2,$3,$4,$5)
		`,
		agentID,
		eventType,
		title,
		description,
		result,
	)

	return err
}
