package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/everest/bheri/modules/webhook/dto"
	"github.com/everest/bheri/modules/webhook/repository"
)

type WebhookService struct {
	repository *repository.WebhookRepository
}

func NewWebhookService(
	repository *repository.WebhookRepository,
) *WebhookService {

	return &WebhookService{
		repository: repository,
	}
}

func (s *WebhookService) Process(
	ctx context.Context,
	req dto.WebhookRequest,
) (*dto.WebhookResponse, error) {

	if req.EventID == "" {
		return nil, errors.New("event_id is required")
	}

	if req.EventType == "" {
		return nil, errors.New("event_type is required")
	}

	if req.Network == "" {
		return nil, errors.New("network is required")
	}

	exists, err := s.repository.Exists(
		ctx,
		req.EventID,
	)

	if err != nil {
		return nil, err
	}

	// Idempotency.
	if exists {
		return &dto.WebhookResponse{
			EventID:   req.EventID,
			EventType: req.EventType,
			Status:    "already_processed",
			Message:   "webhook event already exists",
		}, nil
	}

	id, err := s.repository.Create(
		ctx,
		req.EventID,
		req.EventType,
		req.Network,
		req.Signature,
		req.Data,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create webhook event: %w",
			err,
		)
	}

	// ---------------------------------------------------------
	// BHERI EVENT PROCESSING
	// ---------------------------------------------------------

	switch req.EventType {

	case "payment.confirmed":
		err = s.processPaymentConfirmed(
			ctx,
			req,
		)

	case "payment.failed":
		err = s.processPaymentFailed(
			ctx,
			req,
		)

	case "transaction.confirmed":
		err = s.processTransactionConfirmed(
			ctx,
			req,
		)

	default:
		// Unknown events are intentionally stored.
		err = nil
	}

	if err != nil {

		_ = s.repository.MarkFailed(
			ctx,
			id,
			err.Error(),
		)

		return nil, err
	}

	if err := s.repository.MarkProcessed(
		ctx,
		id,
	); err != nil {
		return nil, err
	}

	return &dto.WebhookResponse{
		ID:        id,
		EventID:   req.EventID,
		EventType: req.EventType,
		Status:    "processed",
		Message:   "webhook processed successfully",
	}, nil
}

func (s *WebhookService) processPaymentConfirmed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	// Later this can call the payment service.
	//
	// Example:
	//
	// paymentID := req.Data["payment_id"]
	// txHash := req.Data["tx_hash"]
	//
	// paymentService.Confirm(...)

	return nil
}

func (s *WebhookService) processPaymentFailed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	// Later:
	// paymentService.Fail(...)

	return nil
}

func (s *WebhookService) processTransactionConfirmed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	// Later:
	// settlementService.Verify(...)
	// transactionService.Confirm(...)

	return nil
}

func (s *WebhookService) Get(
	ctx context.Context,
	id string,
) (map[string]interface{}, error) {

	return s.repository.Get(
		ctx,
		id,
	)
}
