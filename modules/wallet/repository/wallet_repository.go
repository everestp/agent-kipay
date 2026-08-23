// modules/wallet/repository/wallet_repository.go

package repository

import (
	"context"

	"github.com/everest/bheri/models"
	"github.com/jackc/pgx/v5"
)

type WalletRepository interface {
	Create(
		ctx context.Context,
		userID string,
		name string,
		address string,
		network string,
	) (*models.Wallet, error)

	FindByID(
		ctx context.Context,
		id string,
		userID string,
	) (*models.Wallet, error)

	FindByUserID(
		ctx context.Context,
		userID string,
	) ([]models.Wallet, error)
}

type walletRepository struct {
	db *pgx.Conn
}

func NewWalletRepository(
	db *pgx.Conn,
) WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) Create(
	ctx context.Context,
	userID string,
	name string,
	address string,
	network string,
) (*models.Wallet, error) {

	wallet := &models.Wallet{}

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO wallets (
			user_id,
			name,
			address,
			network
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			user_id,
			name,
			address,
			network,
			status,
			created_at,
			updated_at
		`,
		userID,
		name,
		address,
		network,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Address,
		&wallet.Network,
		&wallet.Status,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *walletRepository) FindByID(
	ctx context.Context,
	id string,
	userID string,
) (*models.Wallet, error) {

	wallet := &models.Wallet{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			name,
			address,
			network,
			status,
			created_at,
			updated_at
		FROM wallets
		WHERE id = $1
		AND user_id = $2
		`,
		id,
		userID,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Address,
		&wallet.Network,
		&wallet.Status,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *walletRepository) FindByUserID(
	ctx context.Context,
	userID string,
) ([]models.Wallet, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			user_id,
			name,
			address,
			network,
			status,
			created_at,
			updated_at
		FROM wallets
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wallets := make([]models.Wallet, 0)

	for rows.Next() {

		var wallet models.Wallet

		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Name,
			&wallet.Address,
			&wallet.Network,
			&wallet.Status,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	return wallets, rows.Err()
}
