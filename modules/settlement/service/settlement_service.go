package service

import (
	"context"
	"errors"

	blockchain "github.com/everest/bheri/modules/blockchain/service"
	"github.com/everest/bheri/modules/settlement/repository"
)

type SettlementService interface {

	Settle(
		ctx context.Context,
		paymentID string,
	) error

	Verify(
		ctx context.Context,
		paymentID string,
	) error
}

type settlementService struct {
	repository repository.SettlementRepository
	blockchain blockchain.BlockchainService
}

func NewSettlementService(
	repository repository.SettlementRepository,
	blockchain blockchain.BlockchainService,
) SettlementService {
	return &settlementService{
		repository: repository,
		blockchain: blockchain,
	}
}

func (s *settlementService) Settle(
	ctx context.Context,
	paymentID string,
) error {

	if paymentID == "" {
		return errors.New("payment id is required")
	}

	settlement, err := s.repository.Create(
		ctx,
		paymentID,
	)

	if err != nil {
		return err
	}

	result, err := s.blockchain.Transfer(
		ctx,
		// build transfer from settlement
		nil,
	)

	if err != nil {
		return err
	}

	return s.repository.UpdateSubmitted(
		ctx,
		settlement.ID,
		result.TxHash,
	)
}

func (s *settlementService) Verify(
	ctx context.Context,
	paymentID string,
) error {

	settlement, err := s.repository.FindByPaymentID(
		ctx,
		paymentID,
	)

	if err != nil {
		return err
	}

	if settlement.TxHash == nil {
		return errors.New("transaction not submitted")
	}

	result, err := s.blockchain.VerifyTransaction(
		ctx,
		dto.VerifyTransactionRequest{
			TxHash:  *settlement.TxHash,
			Network: settlement.Network,
		},
	)

	if err != nil {
		return err
	}

	if !result.Confirmed {
		return errors.New("transaction not confirmed")
	}

	return s.repository.UpdateConfirmed(
		ctx,
		settlement.ID,
		int64(result.BlockNumber),
		result.Confirmations,
	)
}
