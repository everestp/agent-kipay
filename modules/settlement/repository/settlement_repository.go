package repository

import (
	"context"

	"github.com/everest/bheri/modules/settlement/dto"
	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type SettlementRepository interface {

	Create(
		ctx context.Context,
		paymentID string,
	) (*models.Settlement, error)

	FindByPaymentID(
		ctx context.Context,
		paymentID string,
	) (*models.Settlement, error)

	UpdateSubmitted(
		ctx context.Context,
		id string,
		txHash string,
	) error

	UpdateConfirmed(
		ctx context.Context,
		id string,
		blockNumber int64,
		confirmations int,
	) error

	UpdateFailed(
		ctx context.Context,
		id string,
		message string,
	) error
}

type settlementRepository struct {
	db *pgx.Conn
}

func NewSettlementRepository(
	db *pgx.Conn,
) SettlementRepository {
	return &settlementRepository{db: db}
}

func (r *settlementRepository) Create(
	ctx context.Context,
	paymentID string,
) (*models.Settlement, error) {

	var settlement models.Settlement

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO settlements (
			payment_id,
			network,
			amount,
			asset,
			sender_address,
			receiver_address
		)
		SELECT
			p.id,
			p.network,
			p.amount,
			p.asset,
			aw.address,
			w.address
		FROM payments p
		JOIN agent_wallets aw
			ON aw.agent_id = p.agent_id
		JOIN wallets w
			ON w.id = p.wallet_id
		WHERE p.id = $1
		RETURNING
			id,
			payment_id,
			tx_hash,
			network,
			amount,
			asset,
			sender_address,
			receiver_address,
			status,
			confirmations,
			block_number,
			created_at
		`,
		paymentID,
	).Scan(
		&settlement.ID,
		&settlement.PaymentID,
		&settlement.TxHash,
		&settlement.Network,
		&settlement.Amount,
		&settlement.Asset,
		&settlement.SenderAddress,
		&settlement.ReceiverAddress,
		&settlement.Status,
		&settlement.Confirmations,
		&settlement.BlockNumber,
		&settlement.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &settlement, nil
}
