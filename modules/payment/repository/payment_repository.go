package repository

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(
		ctx context.Context,
		payment *models.Payment,
	) (*models.Payment, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*models.Payment, error)

	FindByIdempotencyKey(
		ctx context.Context,
		key string,
	) (*models.Payment, error)

	UpdateApproved(
		ctx context.Context,
		id string,
		decision string,
		reason string,
	) error

	UpdateSubmitted(
		ctx context.Context,
		id string,
		txHash string,
	) error

	UpdateCompleted(
		ctx context.Context,
		id string,
	) error

	UpdateFailed(
		ctx context.Context,
		id string,
	) error
}

type paymentRepository struct {
	db *pgxpool.Pool
}
var ErrPaymentNotFound = errors.New("payment not found")
func NewPaymentRepository(
	db *pgxpool.Pool,
) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(
	ctx context.Context,
	payment *models.Payment,
) (*models.Payment, error) {

	var result models.Payment

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO payments (
			agent_id,
			session_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			idempotency_key
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,'pending',$8
		)
		RETURNING
			id,
			agent_id,
			session_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			idempotency_key,
			created_at
		`,
		payment.AgentID,
		payment.SessionID,
		payment.ServiceID,
		payment.Amount,
		payment.Asset,
		payment.Network,
		"x402",
		payment.IdempotencyKey,
	).Scan(
		&result.ID,
		&result.AgentID,
		&result.SessionID,
		&result.ServiceID,
		&result.Amount,
		&result.Asset,
		&result.Network,
		&result.Protocol,
		&result.Status,
		&result.IdempotencyKey,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *paymentRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Payment, error) {

	var payment models.Payment

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			agent_id,
			session_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			COALESCE(policy_decision, ''),
			COALESCE(policy_reason, ''),
			idempotency_key,
			COALESCE(payment_nonce, ''),
			tx_hash,
			created_at,
			approved_at,
			submitted_at,
			completed_at,
			failed_at
		FROM payments
		WHERE id = $1
		`,
		id,
	).Scan(
		&payment.ID,
		&payment.AgentID,
		&payment.SessionID,
		&payment.ServiceID,
		&payment.Amount,
		&payment.Asset,
		&payment.Network,
		&payment.Protocol,
		&payment.Status,
		&payment.PolicyDecision,
		&payment.PolicyReason,
		&payment.IdempotencyKey,
		&payment.PaymentNonce,
		&payment.TxHash,
		&payment.CreatedAt,
		&payment.ApprovedAt,
		&payment.SubmittedAt,
		&payment.CompletedAt,
		&payment.FailedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) FindByIdempotencyKey(
	ctx context.Context,
	key string,
) (*models.Payment, error) {

	var payment models.Payment

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			agent_id,
			session_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			COALESCE(policy_decision, ''),
			COALESCE(policy_reason, ''),
			idempotency_key,
			COALESCE(payment_nonce, ''),
			tx_hash,
			created_at,
			approved_at,
			submitted_at,
			completed_at,
			failed_at
		FROM payments
		WHERE idempotency_key = $1
		`,
		key,
	).Scan(
		&payment.ID,
		&payment.AgentID,
		&payment.SessionID,
		&payment.ServiceID,
		&payment.Amount,
		&payment.Asset,
		&payment.Network,
		&payment.Protocol,
		&payment.Status,
		&payment.PolicyDecision,
		&payment.PolicyReason,
		&payment.IdempotencyKey,
		&payment.PaymentNonce,
		&payment.TxHash,
		&payment.CreatedAt,
		&payment.ApprovedAt,
		&payment.SubmittedAt,
		&payment.CompletedAt,
		&payment.FailedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) UpdateApproved(
	ctx context.Context,
	id string,
	decision string,
	reason string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			status = 'approved',
			policy_decision = $1,
			policy_reason = $2,
			approved_at = NOW()
		WHERE id = $3
		`,
		decision,
		reason,
		id,
	)

	return err
}

func (r *paymentRepository) UpdateSubmitted(
	ctx context.Context,
	id string,
	txHash string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			status = 'submitted',
			tx_hash = $1,
			submitted_at = NOW()
		WHERE id = $2
		`,
		txHash,
		id,
	)

	return err
}

func (r *paymentRepository) UpdateCompleted(
	ctx context.Context,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			status = 'completed',
			completed_at = NOW()
		WHERE id = $1
		`,
		id,
	)

	return err
}

func (r *paymentRepository) UpdateFailed(
	ctx context.Context,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE payments
		SET
			status = 'failed',
			failed_at = NOW()
		WHERE id = $1
		`,
		id,
	)

	return err
}
