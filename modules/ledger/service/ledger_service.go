// modules/ledger/service/ledger_service.go

package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type LedgerService interface {

	RecordSettlement(
		ctx context.Context,
		paymentID string,
		settlementID string,
		agentAccountID string,
		settlementAccountID string,
		amount float64,
	) error
}

type ledgerService struct {
	db *pgx.Conn
}

func NewLedgerService(
	db *pgx.Conn,
) LedgerService {
	return &ledgerService{
		db: db,
	}
}

func (s *ledgerService) RecordSettlement(
	ctx context.Context,
	paymentID string,
	settlementID string,
	agentAccountID string,
	settlementAccountID string,
	amount float64,
) error {

	if amount <= 0 {
		return errors.New("invalid ledger amount")
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	var agentBalance float64

	err = tx.QueryRow(
		ctx,
		`
		SELECT balance
		FROM ledger_accounts
		WHERE id = $1
		FOR UPDATE
		`,
		agentAccountID,
	).Scan(&agentBalance)

	if err != nil {
		return err
	}

	if agentBalance < amount {
		return errors.New("insufficient ledger balance")
	}

	newAgentBalance := agentBalance - amount

	_, err = tx.Exec(
		ctx,
		`
		UPDATE ledger_accounts
		SET
			balance = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		newAgentBalance,
		agentAccountID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO ledger_entries (
			payment_id,
			settlement_id,
			account_id,
			entry_type,
			debit,
			credit,
			balance_after,
			reference
		)
		VALUES (
			$1,$2,$3,'payment',$4,0,$5,$6
		)
		`,
		paymentID,
		settlementID,
		agentAccountID,
		amount,
		newAgentBalance,
		"x402 payment",
	)

	if err != nil {
		return err
	}

	var settlementBalance float64

	err = tx.QueryRow(
		ctx,
		`
		SELECT balance
		FROM ledger_accounts
		WHERE id = $1
		FOR UPDATE
		`,
		settlementAccountID,
	).Scan(&settlementBalance)

	if err != nil {
		return err
	}

	newSettlementBalance :=
		settlementBalance + amount

	_, err = tx.Exec(
		ctx,
		`
		UPDATE ledger_accounts
		SET
			balance = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		newSettlementBalance,
		settlementAccountID,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO ledger_entries (
			payment_id,
			settlement_id,
			account_id,
			entry_type,
			debit,
			credit,
			balance_after,
			reference
		)
		VALUES (
			$1,$2,$3,'settlement',0,$4,$5,$6
		)
		`,
		paymentID,
		settlementID,
		settlementAccountID,
		amount,
		newSettlementBalance,
		"solana settlement",
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
