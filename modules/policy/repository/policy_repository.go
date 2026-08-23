package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type PolicyRepository interface {
	Create(
		ctx context.Context,
		agentID string,
		request PolicyCreateData,
	) (*models.Policy, error)

	FindByAgentID(
		ctx context.Context,
		agentID string,
	) (*models.Policy, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*models.Policy, error)

	UpdateStatus(
		ctx context.Context,
		id string,
		status string,
	) error
}

type PolicyCreateData struct {
	DailyLimit                  float64
	PerTransactionLimit         float64
	WeeklyLimit                 float64
	RequireApprovalAbove        float64
	ExpirationDays              int
	RequireApprovalNewMerchants bool
	RequireApprovalNewAssets    bool
	BlockUnknownAPIs            bool
	AutoPayments                bool
}

type policyRepository struct {
	db *pgx.Conn
}

func NewPolicyRepository(
	db *pgx.Conn,
) PolicyRepository {
	return &policyRepository{
		db: db,
	}
}

func (r *policyRepository) Create(
	ctx context.Context,
	agentID string,
	data PolicyCreateData,
) (*models.Policy, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var policy models.Policy

	err = tx.QueryRow(
		ctx,
		`
		INSERT INTO policies (
			agent_id,
			daily_limit,
			per_transaction_limit,
			weekly_limit,
			require_approval_above,
			expiration_days,
			require_approval_new_merchants,
			require_approval_new_assets,
			block_unknown_apis,
			auto_payments,
			status,
			expires_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			'active',
			NOW() + ($6 * INTERVAL '1 day')
		)
		RETURNING
			id,
			agent_id,
			daily_limit,
			per_transaction_limit,
			weekly_limit,
			require_approval_above,
			expiration_days,
			require_approval_new_merchants,
			require_approval_new_assets,
			block_unknown_apis,
			auto_payments,
			status,
			expires_at,
			created_at,
			updated_at
		`,
		agentID,
		data.DailyLimit,
		data.PerTransactionLimit,
		data.WeeklyLimit,
		data.RequireApprovalAbove,
		data.ExpirationDays,
		data.RequireApprovalNewMerchants,
		data.RequireApprovalNewAssets,
		data.BlockUnknownAPIs,
		data.AutoPayments,
	).Scan(
		&policy.ID,
		&policy.AgentID,
		&policy.DailyLimit,
		&policy.PerTransactionLimit,
		&policy.WeeklyLimit,
		&policy.RequireApprovalAbove,
		&policy.ExpirationDays,
		&policy.RequireApprovalNewMerchants,
		&policy.RequireApprovalNewAssets,
		&policy.BlockUnknownAPIs,
		&policy.AutoPayments,
		&policy.Status,
		&policy.ExpiresAt,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	for _, asset := range []string{"USDC"} {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO policy_allowed_assets
				(policy_id, asset)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			`,
			policy.ID,
			asset,
		)

		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &policy, nil
}

func (r *policyRepository) FindByAgentID(
	ctx context.Context,
	agentID string,
) (*models.Policy, error) {

	var policy models.Policy

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			agent_id,
			daily_limit,
			per_transaction_limit,
			weekly_limit,
			require_approval_above,
			expiration_days,
			require_approval_new_merchants,
			require_approval_new_assets,
			block_unknown_apis,
			auto_payments,
			status,
			expires_at,
			created_at,
			updated_at
		FROM policies
		WHERE agent_id = $1
		AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
		`,
		agentID,
	).Scan(
		&policy.ID,
		&policy.AgentID,
		&policy.DailyLimit,
		&policy.PerTransactionLimit,
		&policy.WeeklyLimit,
		&policy.RequireApprovalAbove,
		&policy.ExpirationDays,
		&policy.RequireApprovalNewMerchants,
		&policy.RequireApprovalNewAssets,
		&policy.BlockUnknownAPIs,
		&policy.AutoPayments,
		&policy.Status,
		&policy.ExpiresAt,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (r *policyRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Policy, error) {

	var policy models.Policy

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			agent_id,
			daily_limit,
			per_transaction_limit,
			weekly_limit,
			require_approval_above,
			expiration_days,
			require_approval_new_merchants,
			require_approval_new_assets,
			block_unknown_apis,
			auto_payments,
			status,
			expires_at,
			created_at,
			updated_at
		FROM policies
		WHERE id = $1
		`,
		id,
	).Scan(
		&policy.ID,
		&policy.AgentID,
		&policy.DailyLimit,
		&policy.PerTransactionLimit,
		&policy.WeeklyLimit,
		&policy.RequireApprovalAbove,
		&policy.ExpirationDays,
		&policy.RequireApprovalNewMerchants,
		&policy.RequireApprovalNewAssets,
		&policy.BlockUnknownAPIs,
		&policy.AutoPayments,
		&policy.Status,
		&policy.ExpiresAt,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (r *policyRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE policies
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		status,
		id,
	)

	return err
}

func placeholders(length int) string {

	values := make([]string, length)

	for i := range values {
		values[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(values, ",")
}
