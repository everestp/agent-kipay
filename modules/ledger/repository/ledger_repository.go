// modules/ledger/repository/ledger_repository.go

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type LedgerRepository interface {

	CreateEntry(
		ctx context.Context,
		tx pgx.Tx,
		paymentID string,
		settlementID string,
		accountID string,
		entryType string,
		debit float64,
		credit float64,
		reference string,
	) error

	GetBalance(
		ctx context.Context,
		accountID string,
	) (float64, error)
}
