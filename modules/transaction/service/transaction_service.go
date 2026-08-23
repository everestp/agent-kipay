package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/modules/transaction/dto"
	"github.com/everest/bheri/modules/transaction/repository"
)

type TransactionService struct {
	repository *repository.TransactionRepository
}

func NewTransactionService(
	repository *repository.TransactionRepository,
) *TransactionService {

	return &TransactionService{
		repository: repository,
	}
}

func (s *TransactionService) Create(
	ctx context.Context,
	req dto.CreateTransactionRequest,
) (*dto.TransactionResponse, error) {

	if req.PaymentID == "" {
		return nil, errors.New("payment_id is required")
	}

	if req.AgentID == "" {
		return nil, errors.New("agent_id is required")
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if req.Asset == "" {
		return nil, errors.New("asset is required")
	}

	if req.Network == "" {
		return nil, errors.New("network is required")
	}

	if req.Protocol == "" {
		return nil, errors.New("protocol is required")
	}

	if req.SenderAddress == "" {
		return nil, errors.New("sender_address is required")
	}

	if req.ReceiverAddress == "" {
		return nil, errors.New("receiver_address is required")
	}

	return s.repository.Create(
		ctx,
		req,
	)
}

func (s *TransactionService) Get(
	ctx context.Context,
	id string,
) (*dto.TransactionResponse, error) {

	if id == "" {
		return nil, errors.New("transaction id is required")
	}

	return s.repository.GetByID(
		ctx,
		id,
	)
}

func (s *TransactionService) List(
	ctx context.Context,
	agentID string,
) ([]dto.TransactionResponse, error) {

	if agentID == "" {
		return nil, errors.New("agent_id is required")
	}

	return s.repository.List(
		ctx,
		agentID,
	)
}

func (s *TransactionService) UpdateStatus(
	ctx context.Context,
	id string,
	req dto.UpdateTransactionRequest,
) error {

	if id == "" {
		return errors.New("transaction id is required")
	}

	if req.Status == "" {
		return errors.New("status is required")
	}

	return s.repository.UpdateStatus(
		ctx,
		id,
		req,
	)
}
