package service

import (
	"context"

	"github.com/everest/bheri/modules/ledger/dto"
	"github.com/everest/bheri/modules/ledger/repository"
)

type LedgerService struct {
	repository *repository.LedgerRepository
}

func NewLedgerService(
	repository *repository.LedgerRepository,
) *LedgerService {

	return &LedgerService{
		repository: repository,
	}
}

// ============================================================
// ACCOUNT
// ============================================================

func (s *LedgerService) CreateAccount(
	ctx context.Context,
	req dto.CreateLedgerAccountRequest,
) (*dto.LedgerAccountResponse, error) {

	return s.repository.CreateAccount(ctx, req)
}

func (s *LedgerService) GetAccount(
	ctx context.Context,
	id string,
) (*dto.LedgerAccountResponse, error) {

	return s.repository.GetAccount(ctx, id)
}

func (s *LedgerService) ListAccounts(
	ctx context.Context,
) ([]dto.LedgerAccountResponse, error) {

	return s.repository.ListAccounts(ctx)
}

// ============================================================
// TRANSACTION
// ============================================================

func (s *LedgerService) CreateTransaction(
	ctx context.Context,
	req dto.CreateLedgerTransactionRequest,
) (*dto.LedgerTransactionResponse, error) {

	return s.repository.CreateTransaction(ctx, req)
}

func (s *LedgerService) GetTransaction(
	ctx context.Context,
	id string,
) (*dto.LedgerTransactionResponse, error) {

	return s.repository.GetTransaction(ctx, id)
}

// ============================================================
// ENTRIES
// ============================================================

func (s *LedgerService) ListAccountEntries(
	ctx context.Context,
	accountID string,
) ([]dto.LedgerEntryResponse, error) {

	return s.repository.ListAccountEntries(ctx, accountID)
}
