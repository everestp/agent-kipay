package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// ============================================================
// PROCESS WEBHOOK
// ============================================================

func (s *WebhookService) Process(
	ctx context.Context,
	req dto.WebhookRequest,
) (*dto.WebhookResponse, error) {

	if s.repository == nil {
		return nil, errors.New("webhook repository is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// VALIDATION
	// ----------------------------------------------------------

	eventID := strings.TrimSpace(req.EventID)
	if eventID == "" {
		return nil, errors.New("event_id is required")
	}

	eventType := strings.TrimSpace(req.EventType)
	if eventType == "" {
		return nil, errors.New("event_type is required")
	}

	network := strings.TrimSpace(req.Network)
	if network == "" {
		return nil, errors.New("network is required")
	}

	// ----------------------------------------------------------
	// IDEMPOTENCY CHECK
	// ----------------------------------------------------------

	exists, err := s.repository.Exists(
		ctx,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check webhook event: %w",
			err,
		)
	}

	if exists {
		return &dto.WebhookResponse{
			EventID:   eventID,
			EventType: eventType,
			Status:    "already_processed",
			Message:   "webhook event already exists",
		}, nil
	}

	// ----------------------------------------------------------
	// CREATE WEBHOOK EVENT
	// ----------------------------------------------------------

	id, err := s.repository.Create(
		ctx,
		eventID,
		eventType,
		network,
		req.Signature,
		req.Data,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create webhook event: %w",
			err,
		)
	}

	// ----------------------------------------------------------
	// PROCESS EVENT
	// ----------------------------------------------------------

	switch eventType {

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
		// Unknown events are stored but not rejected.
		err = nil
	}

	// ----------------------------------------------------------
	// HANDLE PROCESSING ERROR
	// ----------------------------------------------------------

	if err != nil {

		markErr := s.repository.MarkFailed(
			ctx,
			id,
			err.Error(),
		)

		if markErr != nil {
			return nil, errors.Join(
				err,
				fmt.Errorf(
					"mark webhook failed: %w",
					markErr,
				),
			)
		}

		return nil, err
	}

	// ----------------------------------------------------------
	// MARK PROCESSED
	// ----------------------------------------------------------

	if err := s.repository.MarkProcessed(
		ctx,
		id,
	); err != nil {
		return nil, fmt.Errorf(
			"mark webhook processed: %w",
			err,
		)
	}

	// ----------------------------------------------------------
	// RESPONSE
	// ----------------------------------------------------------

	return &dto.WebhookResponse{
		ID:        id,
		EventID:   eventID,
		EventType: eventType,
		Status:    "processed",
		Message:   "webhook processed successfully",
	}, nil
}

// ============================================================
// PAYMENT CONFIRMED
// ============================================================

func (s *WebhookService) processPaymentConfirmed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	/*
		Expected webhook data can contain:

		payment_id
		tx_hash
		amount
		asset
		receiver
		block_number

		Example future implementation:

		paymentID := req.Data["payment_id"]
		txHash := req.Data["tx_hash"]

		paymentService.Confirm(...)
	*/

	return nil
}

// ============================================================
// PAYMENT FAILED
// ============================================================

func (s *WebhookService) processPaymentFailed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	/*
		Expected webhook data:

		payment_id
		tx_hash
		reason

		Future:

		paymentService.Fail(...)
	*/

	return nil
}

// ============================================================
// TRANSACTION CONFIRMED
// ============================================================

func (s *WebhookService) processTransactionConfirmed(
	ctx context.Context,
	req dto.WebhookRequest,
) error {

	if err := ctx.Err(); err != nil {
		return err
	}

	/*
		Expected webhook data:

		transaction_id
		payment_id
		tx_hash
		amount
		asset
		receiver
		block_number

		Future:

		transactionService.Confirm(...)

		settlementService.Verify(...)
	*/

	return nil
}

// ============================================================
// GET WEBHOOK
// ============================================================

func (s *WebhookService) Get(
	ctx context.Context,
	id string,
) (map[string]interface{}, error) {

	if s.repository == nil {
		return nil, errors.New("webhook repository is nil")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, errors.New("webhook id is required")
	}

	result, err := s.repository.Get(
		ctx,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get webhook: %w",
			err,
		)
	}

	return result, nil
}
