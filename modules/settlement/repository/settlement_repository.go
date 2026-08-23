package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Settlement struct {
	ID          string
	PaymentID   string
	TxHash      string
	Network     string
	Status      string
	Amount      float64
	Asset       string
	Receiver    string
	BlockNumber uint64
	Message     string
	CreatedAt   string
	ConfirmedAt *string
	UserID      string
}

type SettlementRepository interface {
	Create(
		ctx context.Context,
		settlement Settlement,
	) (Settlement, error)

	GetByID(
		ctx context.Context,
		userID string,
		id string,
	) (Settlement, error)

	List(
		ctx context.Context,
		userID string,
	) ([]Settlement, error)

	UpdateStatus(
		ctx context.Context,
		id string,
		status string,
		amount float64,
		blockNumber uint64,
		message string,
	) error
}

type settlementRepository struct {
	db *pgxpool.Pool
}

func NewSettlementRepository(
	db *pgxpool.Pool,
) SettlementRepository {
	return &settlementRepository{
		db: db,
	}
}

// ============================================================
// CREATE
// ============================================================

func (r *settlementRepository) Create(
	ctx context.Context,
	s Settlement,
) (Settlement, error) {

	const query = `
		INSERT INTO settlements (
			user_id,
			payment_id,
			tx_hash,
			network,
			status,
			amount,
			asset,
			receiver,
			block_number,
			message
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
			$10
		)
		RETURNING
			id,
			user_id,
			payment_id,
			tx_hash,
			network,
			status,
			amount,
			asset,
			receiver,
			block_number,
			message,
			created_at,
			confirmed_at
	`

	var result Settlement

	err := r.db.QueryRow(
		ctx,
		query,
		s.UserID,
		s.PaymentID,
		s.TxHash,
		s.Network,
		s.Status,
		s.Amount,
		s.Asset,
		s.Receiver,
		s.BlockNumber,
		s.Message,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.PaymentID,
		&result.TxHash,
		&result.Network,
		&result.Status,
		&result.Amount,
		&result.Asset,
		&result.Receiver,
		&result.BlockNumber,
		&result.Message,
		&result.CreatedAt,
		&result.ConfirmedAt,
	)

	if err != nil {
		return Settlement{}, fmt.Errorf(
			"create settlement: %w",
			err,
		)
	}

	return result, nil
}

// ============================================================
// GET BY ID
// ============================================================

func (r *settlementRepository) GetByID(
	ctx context.Context,
	userID string,
	id string,
) (Settlement, error) {

	const query = `
		SELECT
			id,
			user_id,
			payment_id,
			tx_hash,
			network,
			status,
			amount,
			asset,
			receiver,
			block_number,
			message,
			created_at,
			confirmed_at
		FROM settlements
		WHERE id = $1
		AND user_id = $2
	`

	var settlement Settlement

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&settlement.ID,
		&settlement.UserID,
		&settlement.PaymentID,
		&settlement.TxHash,
		&settlement.Network,
		&settlement.Status,
		&settlement.Amount,
		&settlement.Asset,
		&settlement.Receiver,
		&settlement.BlockNumber,
		&settlement.Message,
		&settlement.CreatedAt,
		&settlement.ConfirmedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return Settlement{}, pgx.ErrNoRows
		}

		return Settlement{}, fmt.Errorf(
			"get settlement: %w",
			err,
		)
	}

	return settlement, nil
}

// ============================================================
// LIST
// ============================================================

func (r *settlementRepository) List(
	ctx context.Context,
	userID string,
) ([]Settlement, error) {

	const query = `
		SELECT
			id,
			user_id,
			payment_id,
			tx_hash,
			network,
			status,
			amount,
			asset,
			receiver,
			block_number,
			message,
			created_at,
			confirmed_at
		FROM settlements
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"list settlements: %w",
			err,
		)
	}

	defer rows.Close()

	settlements := make(
		[]Settlement,
		0,
	)

	for rows.Next() {

		var settlement Settlement

		err := rows.Scan(
			&settlement.ID,
			&settlement.UserID,
			&settlement.PaymentID,
			&settlement.TxHash,
			&settlement.Network,
			&settlement.Status,
			&settlement.Amount,
			&settlement.Asset,
			&settlement.Receiver,
			&settlement.BlockNumber,
			&settlement.Message,
			&settlement.CreatedAt,
			&settlement.ConfirmedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan settlement: %w",
				err,
			)
		}

		settlements = append(
			settlements,
			settlement,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate settlements: %w",
			err,
		)
	}

	return settlements, nil
}

// ============================================================
// UPDATE STATUS
// ============================================================

func (r *settlementRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
	amount float64,
	blockNumber uint64,
	message string,
) error {

	const query = `
		UPDATE settlements
		SET
			status = $1,
			amount = $2,
			block_number = $3,
			message = $4,

			confirmed_at =
				CASE
					WHEN $1 = 'confirmed'
						AND confirmed_at IS NULL
					THEN NOW()

					WHEN $1 != 'confirmed'
					THEN confirmed_at

					ELSE confirmed_at
				END

		WHERE id = $5
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		status,
		amount,
		blockNumber,
		message,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"update settlement status: %w",
			err,
		)
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
