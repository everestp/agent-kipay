package repository

import (
	"context"
	
	"fmt"

	"github.com/everest/bheri/modules/transaction/dto"
	"github.com/jackc/pgx/v5"
)

type TransactionRepository struct {
	db *pgx.Conn
}

func NewTransactionRepository(db *pgx.Conn) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	req dto.CreateTransactionRequest,
) (*dto.TransactionResponse, error) {

	query := `
		INSERT INTO transactions (
			payment_id,
			agent_id,
			session_id,
			api_service_id,
			type,
			protocol,
			network,
			asset,
			amount,
			sender_address,
			receiver_address,
			status,
			verification_status,
			settlement_status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			'pending',
			'pending',
			'pending'
		)
		RETURNING
			id,
			payment_id,
			agent_id,
			session_id,
			api_service_id,
			type,
			protocol,
			network,
			asset,
			amount,
			sender_address,
			receiver_address,
			status,
			verification_status,
			settlement_status,
			created_at,
			updated_at
	`

	var tx dto.TransactionResponse

	err := r.db.QueryRow(
		ctx,
		query,
		req.PaymentID,
		req.AgentID,
		req.SessionID,
		req.APIServiceID,
		req.Type,
		req.Protocol,
		req.Network,
		req.Asset,
		req.Amount,
		req.SenderAddress,
		req.ReceiverAddress,
	).Scan(
		&tx.ID,
		&tx.PaymentID,
		&tx.AgentID,
		&tx.SessionID,
		&tx.APIServiceID,
		&tx.Type,
		&tx.Protocol,
		&tx.Network,
		&tx.Asset,
		&tx.Amount,
		&tx.SenderAddress,
		&tx.ReceiverAddress,
		&tx.Status,
		&tx.VerificationStatus,
		&tx.SettlementStatus,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create transaction: %w",
			err,
		)
	}

	return &tx, nil
}

func (r *TransactionRepository) GetByID(
	ctx context.Context,
	id string,
) (*dto.TransactionResponse, error) {

	query := `
		SELECT
			id,
			payment_id,
			agent_id,
			session_id,
			api_service_id,
			type,
			protocol,
			network,
			asset,
			amount,
			sender_address,
			receiver_address,
			status,
			tx_hash,
			block_number,
			verification_status,
			settlement_status,
			created_at,
			updated_at
		FROM transactions
		WHERE id = $1
	`

	var tx dto.TransactionResponse

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&tx.ID,
		&tx.PaymentID,
		&tx.AgentID,
		&tx.SessionID,
		&tx.APIServiceID,
		&tx.Type,
		&tx.Protocol,
		&tx.Network,
		&tx.Asset,
		&tx.Amount,
		&tx.SenderAddress,
		&tx.ReceiverAddress,
		&tx.Status,
		&tx.TxHash,
		&tx.BlockNumber,
		&tx.VerificationStatus,
		&tx.SettlementStatus,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get transaction: %w",
			err,
		)
	}

	return &tx, nil
}

func (r *TransactionRepository) List(
	ctx context.Context,
	agentID string,
) ([]dto.TransactionResponse, error) {

	query := `
		SELECT
			id,
			payment_id,
			agent_id,
			session_id,
			api_service_id,
			type,
			protocol,
			network,
			asset,
			amount,
			sender_address,
			receiver_address,
			status,
			tx_hash,
			block_number,
			verification_status,
			settlement_status,
			created_at,
			updated_at
		FROM transactions
		WHERE agent_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		agentID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"list transactions: %w",
			err,
		)
	}

	defer rows.Close()

	transactions := make(
		[]dto.TransactionResponse,
		0,
	)

	for rows.Next() {

		var tx dto.TransactionResponse

		err := rows.Scan(
			&tx.ID,
			&tx.PaymentID,
			&tx.AgentID,
			&tx.SessionID,
			&tx.APIServiceID,
			&tx.Type,
			&tx.Protocol,
			&tx.Network,
			&tx.Asset,
			&tx.Amount,
			&tx.SenderAddress,
			&tx.ReceiverAddress,
			&tx.Status,
			&tx.TxHash,
			&tx.BlockNumber,
			&tx.VerificationStatus,
			&tx.SettlementStatus,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan transaction: %w",
				err,
			)
		}

		transactions = append(
			transactions,
			tx,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateStatus(
	ctx context.Context,
	id string,
	req dto.UpdateTransactionRequest,
) error {
	query := `
		UPDATE transactions
		SET
			status = $2,
			tx_hash = COALESCE($3, tx_hash),
			block_number = COALESCE($4, block_number),
			verification_status =
				CASE
					WHEN $5 = '' THEN verification_status
					ELSE $5
				END,
			settlement_status =
				CASE
					WHEN $6 = '' THEN settlement_status
					ELSE $6
				END,
			updated_at = NOW()
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		id,
		req.Status,
		req.TxHash,
		req.BlockNumber,
		req.VerificationStatus,
		req.SettlementStatus,
	)
	if err != nil {
		return fmt.Errorf(
			"update transaction: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
