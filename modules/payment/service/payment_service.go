package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/payment/dto"
	"github.com/everest/bheri/modules/payment/repository"
)

type PaymentService interface {
	Create(
		ctx context.Context,
		request dto.CreatePaymentRequest,
	) (*dto.PaymentResponse, error)

	Get(
		ctx context.Context,
		id string,
	) (*models.Payment, error)
}

type paymentService struct {
	repository repository.PaymentRepository
	policy      PolicyEngine
	session     SessionService
}

func NewPaymentService(
	repository repository.PaymentRepository,
	policy PolicyEngine,
	session SessionService,
) PaymentService {
	return &paymentService{
		repository: repository,
		policy:      policy,
		session:     session,
	}
}

func (s *paymentService) Create(
	ctx context.Context,
	request dto.CreatePaymentRequest,
) (*dto.PaymentResponse, error) {

	if request.AgentID == "" {
		return nil, errors.New("agent id is required")
	}

	if request.SessionID == "" {
		return nil, errors.New("session id is required")
	}

	if request.Amount <= 0 {
		return nil, errors.New("invalid payment amount")
	}

	if request.IdempotencyKey == "" {
		return nil, errors.New("idempotency key is required")
	}

	existing, err := s.repository.FindByIdempotencyKey(
		ctx,
		request.IdempotencyKey,
	)

	if err == nil {
		return &dto.PaymentResponse{
			ID:             existing.ID,
			Status:         existing.Status,
			Amount:         existing.Amount,
			Asset:          existing.Asset,
			Network:        existing.Network,
			Protocol:       existing.Protocol,
			PolicyDecision: existing.PolicyDecision,
			PolicyReason:   existing.PolicyReason,
			TxHash:         existing.TxHash,
		}, nil
	}

	_, err = s.session.Validate(
		ctx,
		request.SessionID,
	)

	if err != nil {
		return nil, err
	}

	err = s.policy.Authorize(
		ctx,
		request.AgentID,
		request.SessionID,
		request.Amount,
		request.Asset,
		request.Network,
		request.ServiceID,
		"",
	)

	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		AgentID:        request.AgentID,
		SessionID:      &request.SessionID,
		ServiceID:      &request.ServiceID,
		Amount:         request.Amount,
		Asset:          request.Asset,
		Network:        request.Network,
		IdempotencyKey: request.IdempotencyKey,
	}

	created, err := s.repository.Create(
		ctx,
		payment,
	)

	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		ID:             created.ID,
		Status:         created.Status,
		Amount:         created.Amount,
		Asset:          created.Asset,
		Network:        created.Network,
		Protocol:       created.Protocol,
	}, nil
}

func (s *paymentService) Get(
	ctx context.Context,
	id string,
) (*models.Payment, error) {

	return s.repository.FindByID(
		ctx,
		id,
	)
}
