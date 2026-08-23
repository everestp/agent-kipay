package repository

import (
	"context"
	"fmt"

	"github.com/everest/bheri/modules/ledger/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerRepository struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}

// ============================================================
// LEDGER ACCOUNTS
// ============================================================

func (r *LedgerRepository) CreateAccount(
	ctx context.Context,
	req dto.CreateLedgerAccountRequest,
) (*dto.LedgerAccountResponse, error) {

	query := `
		INSERT INTO ledger_accounts (
			wallet_id,
			agent_id,
			name,
			currency
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			wallet_id,
			agent_id,
			name,
			currency,
			created_at
	`

	var account dto.LedgerAccountResponse

	err := r.db.QueryRow(
		ctx,
		query,
		req.WalletID,
		req.AgentID,
		req.Name,
		req.Currency,
	).Scan(
		&account.ID,
		&account.WalletID,
		&account.AgentID,
		&account.Name,
		&account.Currency,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create ledger account: %w",
			err,
		)
	}

	return &account, nil
}

func (r *LedgerRepository) GetAccount(
	ctx context.Context,
	id string,
) (*dto.LedgerAccountResponse, error) {

	query := `
		SELECT
			id,
			wallet_id,
			agent_id,
			name,
			currency,
			created_at
		FROM ledger_accounts
		WHERE id = $1
	`

	var account dto.LedgerAccountResponse

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&account.ID,
		&account.WalletID,
		&account.AgentID,
		&account.Name,
		&account.Currency,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get ledger account: %w",
			err,
		)
	}

	return &account, nil
}

func (r *LedgerRepository) ListAccounts(
	ctx context.Context,
) ([]dto.LedgerAccountResponse, error) {

	query := `
		SELECT
			id,
			wallet_id,
			agent_id,
			name,
			currency,
			created_at
		FROM ledger_accounts
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"list ledger accounts: %w",
			err,
		)
	}

	defer rows.Close()

	accounts := make(
		[]dto.LedgerAccountResponse,
		0,
	)

	for rows.Next() {
		var account dto.LedgerAccountResponse

		err := rows.Scan(
			&account.ID,
			&account.WalletID,
			&account.AgentID,
			&account.Name,
			&account.Currency,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan ledger account: %w",
				err,
			)
		}

		accounts = append(
			accounts,
			account,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate ledger accounts: %w",
			err,
		)
	}

	return accounts, nil
}

// ============================================================
// LEDGER TRANSACTION
// ============================================================

func (r *LedgerRepository) CreateTransaction(
	ctx context.Context,
	req dto.CreateLedgerTransactionRequest,
) (*dto.LedgerTransactionResponse, error) {

	if req.Amount <= 0 {
		return nil, fmt.Errorf(
			"amount must be greater than zero",
		)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"begin ledger transaction: %w",
			err,
		)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// --------------------------------------------------------
	// CREATE LEDGER TRANSACTION
	// --------------------------------------------------------

	insertTransaction := `
		INSERT INTO ledger_transactions (
			payment_id,
			transaction_id,
			reference
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			payment_id,
			transaction_id,
			reference,
			created_at
	`

	var ledgerTx dto.LedgerTransactionResponse

	err = tx.QueryRow(
		ctx,
		insertTransaction,
		req.PaymentID,
		req.TransactionID,
		req.Reference,
	).Scan(
		&ledgerTx.ID,
		&ledgerTx.PaymentID,
		&ledgerTx.TransactionID,
		&ledgerTx.Reference,
		&ledgerTx.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create ledger transaction: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// DEBIT
	// --------------------------------------------------------

	debitQuery := `
		INSERT INTO ledger_entries (
			ledger_transaction_id,
			ledger_account_id,
			entry_type,
			amount,
			asset
		)
		VALUES ($1, $2, 'debit', $3, $4)
		RETURNING
			id,
			ledger_transaction_id,
			ledger_account_id,
			entry_type,
			amount,
			asset,
			created_at
	`

	var debit dto.LedgerEntryResponse

	err = tx.QueryRow(
		ctx,
		debitQuery,
		ledgerTx.ID,
		req.DebitAccountID,
		req.Amount,
		req.Asset,
	).Scan(
		&debit.ID,
		&debit.LedgerTransactionID,
		&debit.LedgerAccountID,
		&debit.EntryType,
		&debit.Amount,
		&debit.Asset,
		&debit.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create debit ledger entry: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// CREDIT
	// --------------------------------------------------------

	creditQuery := `
		INSERT INTO ledger_entries (
			ledger_transaction_id,
			ledger_account_id,
			entry_type,
			amount,
			asset
		)
		VALUES ($1, $2, 'credit', $3, $4)
		RETURNING
			id,
			ledger_transaction_id,
			ledger_account_id,
			entry_type,
			amount,
			asset,
			created_at
	`

	var credit dto.LedgerEntryResponse

	err = tx.QueryRow(
		ctx,
		creditQuery,
		ledgerTx.ID,
		req.CreditAccountID,
		req.Amount,
		req.Asset,
	).Scan(
		&credit.ID,
		&credit.LedgerTransactionID,
		&credit.LedgerAccountID,
		&credit.EntryType,
		&credit.Amount,
		&credit.Asset,
		&credit.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create credit ledger entry: %w",
			err,
		)
	}

	ledgerTx.Entries = []dto.LedgerEntryResponse{
		debit,
		credit,
	}

	// --------------------------------------------------------
	// COMMIT
	// --------------------------------------------------------

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"commit ledger transaction: %w",
			err,
		)
	}

	return &ledgerTx, nil
}

// ============================================================
// GET LEDGER TRANSACTION
// ============================================================

func (r *LedgerRepository) GetTransaction(
	ctx context.Context,
	id string,
) (*dto.LedgerTransactionResponse, error) {

	query := `
		SELECT
			lt.id,
			lt.payment_id,
			lt.transaction_id,
			lt.reference,
			lt.created_at,

			le.id,
			le.ledger_transaction_id,
			le.ledger_account_id,
			le.entry_type,
			le.amount,
			le.asset,
			le.created_at

		FROM ledger_transactions lt

		LEFT JOIN ledger_entries le
			ON le.ledger_transaction_id = lt.id

		WHERE lt.id = $1

		ORDER BY le.created_at ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get ledger transaction: %w",
			err,
		)
	}

	defer rows.Close()

	var result *dto.LedgerTransactionResponse

	for rows.Next() {

		if result == nil {
			result = &dto.LedgerTransactionResponse{}
		}

		var (
			entryID            pgtype.Text
			entryTransactionID pgtype.Text
			entryAccountID     pgtype.Text
			entryType          pgtype.Text
			entryAmount        pgtype.Float8
			entryAsset         pgtype.Text
			entryCreatedAt     pgtype.Text
		)

		err := rows.Scan(
			&result.ID,
			&result.PaymentID,
			&result.TransactionID,
			&result.Reference,
			&result.CreatedAt,

			&entryID,
			&entryTransactionID,
			&entryAccountID,
			&entryType,
			&entryAmount,
			&entryAsset,
			&entryCreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan ledger transaction: %w",
				err,
			)
		}

		if entryID.Valid {

			entry := dto.LedgerEntryResponse{
				ID:                  entryID.String,
				LedgerTransactionID: entryTransactionID.String,
				LedgerAccountID:     entryAccountID.String,
				EntryType:           entryType.String,
				Amount:              entryAmount.Float64,
				Asset:               entryAsset.String,
				CreatedAt:           entryCreatedAt.String,
			}

			result.Entries = append(
				result.Entries,
				entry,
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate ledger transaction: %w",
			err,
		)
	}

	if result == nil {
		return nil, pgx.ErrNoRows
	}

	return result, nil
}

// ============================================================
// ACCOUNT ENTRIES
// ============================================================

func (r *LedgerRepository) ListAccountEntries(
	ctx context.Context,
	accountID string,
) ([]dto.LedgerEntryResponse, error) {

	query := `
		SELECT
			id,
			ledger_transaction_id,
			ledger_account_id,
			entry_type,
			amount,
			asset,
			created_at
		FROM ledger_entries
		WHERE ledger_account_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		accountID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"list account entries: %w",
			err,
		)
	}

	defer rows.Close()

	entries := make(
		[]dto.LedgerEntryResponse,
		0,
	)

	for rows.Next() {

		var entry dto.LedgerEntryResponse

		err := rows.Scan(
			&entry.ID,
			&entry.LedgerTransactionID,
			&entry.LedgerAccountID,
			&entry.EntryType,
			&entry.Amount,
			&entry.Asset,
			&entry.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan ledger entry: %w",
				err,
			)
		}

		entries = append(
			entries,
			entry,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate ledger entries: %w",
			err,
		)
	}

	return entries, nil
}
