package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type WebhookRepository struct {
	db *pgx.Conn
}

func NewWebhookRepository(db *pgx.Conn) *WebhookRepository {
	return &WebhookRepository{
		db: db,
	}
}

// ============================================================
// EXISTS
// ============================================================

func (r *WebhookRepository) Exists(
	ctx context.Context,
	eventID string,
) (bool, error) {

	var exists bool

	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM webhook_events
			WHERE event_id = $1
		)
		`,
		eventID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// ============================================================
// CREATE
// ============================================================

func (r *WebhookRepository) Create(
	ctx context.Context,
	eventID string,
	eventType string,
	network string,
	signature string,
	payload map[string]interface{},
) (string, error) {

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf(
			"marshal webhook payload: %w",
			err,
		)
	}

	var id string

	err = r.db.QueryRow(
		ctx,
		`
		INSERT INTO webhook_events (
			event_id,
			event_type,
			network,
			signature,
			payload,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			'received'
		)
		RETURNING id
		`,
		eventID,
		eventType,
		network,
		signature,
		data,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf(
			"create webhook event: %w",
			err,
		)
	}

	return id, nil
}

// ============================================================
// MARK PROCESSED
// ============================================================

func (r *WebhookRepository) MarkProcessed(
	ctx context.Context,
	id string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE webhook_events
		SET
			status = 'processed',
			processed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		`,
		id,
	)

	if err != nil {
		return fmt.Errorf(
			"mark webhook processed: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// MARK FAILED
// ============================================================

func (r *WebhookRepository) MarkFailed(
	ctx context.Context,
	id string,
	message string,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		UPDATE webhook_events
		SET
			status = 'failed',
			error_message = $2,
			updated_at = NOW()
		WHERE id = $1
		`,
		id,
		message,
	)

	if err != nil {
		return fmt.Errorf(
			"mark webhook failed: %w",
			err,
		)
	}

	return nil
}

// ============================================================
// GET
// ============================================================

func (r *WebhookRepository) Get(
	ctx context.Context,
	id string,
) (map[string]interface{}, error) {

	var (
		eventID      string
		eventType    string
		network      string
		status       string
		payload      []byte
		errorMessage  *string
	)

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			event_id,
			event_type,
			network,
			status,
			payload,
			error_message
		FROM webhook_events
		WHERE id = $1
		`,
		id,
	).Scan(
		&eventID,
		&eventType,
		&network,
		&status,
		&payload,
		&errorMessage,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get webhook event: %w",
			err,
		)
	}

	var data map[string]interface{}

	if len(payload) > 0 {
		if err := json.Unmarshal(
			payload,
			&data,
		); err != nil {
			return nil, fmt.Errorf(
				"invalid webhook payload: %w",
				err,
			)
		}
	} else {
		data = make(map[string]interface{})
	}

	var errorMessageValue string

	if errorMessage != nil {
		errorMessageValue = *errorMessage
	}

	return map[string]interface{}{
		"id":            id,
		"event_id":      eventID,
		"event_type":    eventType,
		"network":       network,
		"status":        status,
		"payload":       data,
		"error_message": errorMessageValue,
	}, nil
}
