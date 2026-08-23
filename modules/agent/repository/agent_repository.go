package repository

import (
	"context"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository interface {
	Create(
		ctx context.Context,
		userID string,
		walletID string,
		name string,
		description string,
		asset string,
		network string,
		color string,
	) (*models.Agent, error)

	FindByID(
		ctx context.Context,
		userID string,
		id string,
	) (*models.Agent, error)

	FindByUserID(
		ctx context.Context,
		userID string,
	) ([]models.Agent, error)

	Update(
		ctx context.Context,
		userID string,
		id string,
		name string,
		description string,
		status string,
		color string,
		autoPayments bool,
	) (*models.Agent, error)

	Delete(
		ctx context.Context,
		userID string,
		id string,
	) error
}

type agentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(
	db *pgxpool.Pool,
) AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(
	ctx context.Context,
	userID string,
	walletID string,
	name string,
	description string,
	asset string,
	network string,
	color string,
) (*models.Agent, error) {

	var agent models.Agent

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO agents (
			user_id,
			wallet_id,
			name,
			description,
			asset,
			network,
			color
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING
			id,
			user_id,
			wallet_id,
			name,
			COALESCE(description, ''),
			status,
			asset,
			network,
			COALESCE(color, ''),
			auto_payments,
			last_active_at,
			created_at,
			updated_at
		`,
		userID,
		walletID,
		name,
		description,
		asset,
		network,
		color,
	).Scan(
		&agent.ID,
		&agent.UserID,
		&agent.WalletID,
		&agent.Name,
		&agent.Description,
		&agent.Status,
		&agent.Asset,
		&agent.Network,
		&agent.Color,
		&agent.AutoPayments,
		&agent.LastActiveAt,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *agentRepository) FindByID(
	ctx context.Context,
	userID string,
	id string,
) (*models.Agent, error) {

	var agent models.Agent

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			wallet_id,
			name,
			COALESCE(description, ''),
			status,
			asset,
			network,
			COALESCE(color, ''),
			auto_payments,
			last_active_at,
			created_at,
			updated_at
		FROM agents
		WHERE id = $1
		AND user_id = $2
		`,
		id,
		userID,
	).Scan(
		&agent.ID,
		&agent.UserID,
		&agent.WalletID,
		&agent.Name,
		&agent.Description,
		&agent.Status,
		&agent.Asset,
		&agent.Network,
		&agent.Color,
		&agent.AutoPayments,
		&agent.LastActiveAt,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *agentRepository) FindByUserID(
	ctx context.Context,
	userID string,
) ([]models.Agent, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			user_id,
			wallet_id,
			name,
			COALESCE(description, ''),
			status,
			asset,
			network,
			COALESCE(color, ''),
			auto_payments,
			last_active_at,
			created_at,
			updated_at
		FROM agents
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	agents := make([]models.Agent, 0)

	for rows.Next() {

		var agent models.Agent

		err := rows.Scan(
			&agent.ID,
			&agent.UserID,
			&agent.WalletID,
			&agent.Name,
			&agent.Description,
			&agent.Status,
			&agent.Asset,
			&agent.Network,
			&agent.Color,
			&agent.AutoPayments,
			&agent.LastActiveAt,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	return agents, rows.Err()
}

func (r *agentRepository) Update(
	ctx context.Context,
	userID string,
	id string,
	name string,
	description string,
	status string,
	color string,
	autoPayments bool,
) (*models.Agent, error) {

	var agent models.Agent

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE agents
		SET
			name = $1,
			description = $2,
			status = $3,
			color = $4,
			auto_payments = $5,
			updated_at = NOW()
		WHERE id = $6
		AND user_id = $7
		RETURNING
			id,
			user_id,
			wallet_id,
			name,
			COALESCE(description, ''),
			status,
			asset,
			network,
			COALESCE(color, ''),
			auto_payments,
			last_active_at,
			created_at,
			updated_at
		`,
		name,
		description,
		status,
		color,
		autoPayments,
		id,
		userID,
	).Scan(
		&agent.ID,
		&agent.UserID,
		&agent.WalletID,
		&agent.Name,
		&agent.Description,
		&agent.Status,
		&agent.Asset,
		&agent.Network,
		&agent.Color,
		&agent.AutoPayments,
		&agent.LastActiveAt,
		&agent.CreatedAt,
		&agent.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *agentRepository) Delete(
	ctx context.Context,
	userID string,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		DELETE FROM agents
		WHERE id = $1
		AND user_id = $2
		`,
		id,
		userID,
	)

	return err
}
