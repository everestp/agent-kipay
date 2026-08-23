// modules/transaction/repository/transaction_repository.go

package repository

import (
	"context"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type TransactionRepository interface {

	CreateFromPayment(
		ctx context.Context,
		paymentID string,
	) (*models.Transaction, error)

	FindByID(
		ctx context.Context,
		id string,
	) (*models.Transaction, error)

	ListByAgent(
		ctx context.Context,
		agentID string,
	) ([]models.Transaction, error)

	Complete(
		ctx context.Context,
		paymentID string,
		txHash string,
	) error

	Fail(
		ctx context.Context,
		paymentID string,
	) error
}

type transactionRepository struct {
	db *pgx.Conn
}

func NewTransactionRepository(
	db *pgx.Conn,
) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateFromPayment(
	ctx context.Context,
	paymentID string,
) (*models.Transaction, error) {

	var transaction models.Transaction

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO transactions (
			payment_id,
			agent_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status
		)
		SELECT
			id,
			agent_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			'pending'
		FROM payments
		WHERE id = $1
		ON CONFLICT (payment_id)
		DO UPDATE SET payment_id = EXCLUDED.payment_id
		RETURNING
			id,
			payment_id,
			settlement_id,
			agent_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			tx_hash,
			created_at,
			completed_at
		`,
		paymentID,
	).Scan(
		&transaction.ID,
		&transaction.PaymentID,
		&transaction.SettlementID,
		&transaction.AgentID,
		&transaction.ServiceID,
		&transaction.Amount,
		&transaction.Asset,
		&transaction.Network,
		&transaction.Protocol,
		&transaction.Status,
		&transaction.TxHash,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByID(
	ctx context.Context,
	id string,
) (*models.Transaction, error) {

	var transaction models.Transaction

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			payment_id,
			settlement_id,
			agent_id,
			service_id,
			amount,
			asset,
			network,
			protocol,
			status,
			tx_hash,
			created_at,
			completed_at
		FROM transactions
		WHERE id = $1
		`,
		id,
	).Scan(
		&transaction.ID,
		&transaction.PaymentID,
		&transaction.SettlementID,
		&transaction.AgentID,
		&transaction.ServiceID,
		&transaction.Amount,
		&transaction.Asset,
		&transaction.Network,
		&transaction.Protocol,
		&transaction.Status,
		&transaction.TxHash,
		&transaction.CreatedAt,
		&transaction.CompletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) Complete(
	ctx context.Context,
	paymentID string,
	txHash string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE transactions
		SET
			status = 'completed',
			tx_hash = $1,
			completed_at = NOW()
		WHERE payment_id = $2
		`,
		txHash,
		paymentID,
	)

	return err
}

func (r *transactionRepository) Fail(
	ctx context.Context,
	paymentID string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE transactions
		SET status = 'failed'
		WHERE payment_id = $1
		`,
		paymentID,
	)

	return err
}
