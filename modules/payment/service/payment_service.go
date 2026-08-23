package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/payment/dto"
	"github.com/everest/bheri/modules/payment/repository"
	policyservice "github.com/everest/bheri/modules/policy/service"
	sessionservice "github.com/everest/bheri/modules/session/service"
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
	policy     policyservice.PolicyEngine
	session    sessionservice.SessionService
}

func NewPaymentService(
	repository repository.PaymentRepository,
	policy policyservice.PolicyEngine,
	session sessionservice.SessionService,
) PaymentService {
	return &paymentService{
		repository: repository,
		policy:     policy,
		session:    session,
	}
}

func (s *paymentService) Create(
	ctx context.Context,
	request dto.CreatePaymentRequest,
) (*dto.PaymentResponse, error) {

	// =========================================================
	// VALIDATION
	// =========================================================

	if request.AgentID == "" {
		return nil, errors.New("agent id is required")
	}

	if request.SessionID == "" {
		return nil, errors.New("session id is required")
	}

	if request.ServiceID == "" {
		return nil, errors.New("service id is required")
	}

	if request.Amount <= 0 {
		return nil, errors.New("invalid payment amount")
	}

	if request.Asset == "" {
		return nil, errors.New("asset is required")
	}

	if request.Network == "" {
		return nil, errors.New("network is required")
	}

	if request.IdempotencyKey == "" {
		return nil, errors.New("idempotency key is required")
	}

	// =========================================================
	// IDEMPOTENCY
	// =========================================================

	existing, err := s.repository.FindByIdempotencyKey(
		ctx,
		request.IdempotencyKey,
	)

	if err == nil && existing != nil {
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

	// =========================================================
	// SESSION VALIDATION
	// =========================================================

	_, err = s.session.Validate(
		ctx,
		request.SessionID,
	)

	if err != nil {
		return nil, err
	}

	// =========================================================
	// POLICY AUTHORIZATION
	// =========================================================

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

	// =========================================================
	// CREATE PAYMENT
	// =========================================================

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

	// =========================================================
	// RESPONSE
	// =========================================================

	return &dto.PaymentResponse{
		ID:       created.ID,
		Status:   created.Status,
		Amount:   created.Amount,
		Asset:    created.Asset,
		Network:  created.Network,
		Protocol: created.Protocol,
	}, nil
}

func (s *paymentService) Get(
	ctx context.Context,
	id string,
) (*models.Payment, error) {

	if id == "" {
		return nil, errors.New("payment id is required")
	}

	return s.repository.FindByID(
		ctx,
		id,
	)
}
