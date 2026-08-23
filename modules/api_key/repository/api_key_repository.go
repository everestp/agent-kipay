package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKey struct {
	ID             string
	UserID         string
	Name           string
	PublishableKey string
	SecretKeyHash  string
	Status         string
	CreatedAt      time.Time
	RevokedAt      *time.Time
}

type APIKeyRepository interface {
	Create(
		ctx context.Context,
		userID string,
		name string,
		publishableKey string,
		secretKeyHash string,
	) (APIKey, error)

	List(
		ctx context.Context,
		userID string,
	) ([]APIKey, error)

	GetByID(
		ctx context.Context,
		userID string,
		id string,
	) (APIKey, error)

	Revoke(
		ctx context.Context,
		userID string,
		id string,
	) error

	GetByPublishableKey(
		ctx context.Context,
		key string,
	) (APIKey, error)

	GetBySecretHash(
		ctx context.Context,
		hash string,
	) (*APIKey, error)
}

type apiKeyRepository struct {
	db *pgxpool.Pool
}

func NewAPIKeyRepository(db *pgxpool.Pool) APIKeyRepository {
	return &apiKeyRepository{
		db: db,
	}
}

// ============================================================
// CREATE
// ============================================================

func (r *apiKeyRepository) Create(
	ctx context.Context,
	userID string,
	name string,
	publishableKey string,
	secretKeyHash string,
) (APIKey, error) {

	var key APIKey

	err := r.db.QueryRow(ctx, `
		INSERT INTO api_keys (
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			'active'
		)
		RETURNING
			id,
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status,
			created_at,
			revoked_at
	`,
		userID,
		name,
		publishableKey,
		secretKeyHash,
	).Scan(
		&key.ID,
		&key.UserID,
		&key.Name,
		&key.PublishableKey,
		&key.SecretKeyHash,
		&key.Status,
		&key.CreatedAt,
		&key.RevokedAt,
	)

	if err != nil {
		return APIKey{}, err
	}

	return key, nil
}

// ============================================================
// LIST
// ============================================================

func (r *apiKeyRepository) List(
	ctx context.Context,
	userID string,
) ([]APIKey, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status,
			created_at,
			revoked_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := make([]APIKey, 0)

	for rows.Next() {
		var key APIKey

		err := rows.Scan(
			&key.ID,
			&key.UserID,
			&key.Name,
			&key.PublishableKey,
			&key.SecretKeyHash,
			&key.Status,
			&key.CreatedAt,
			&key.RevokedAt,
		)

		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

// ============================================================
// GET BY ID
// ============================================================

func (r *apiKeyRepository) GetByID(
	ctx context.Context,
	userID string,
	id string,
) (APIKey, error) {

	var key APIKey

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status,
			created_at,
			revoked_at
		FROM api_keys
		WHERE id = $1
		AND user_id = $2
	`,
		id,
		userID,
	).Scan(
		&key.ID,
		&key.UserID,
		&key.Name,
		&key.PublishableKey,
		&key.SecretKeyHash,
		&key.Status,
		&key.CreatedAt,
		&key.RevokedAt,
	)

	if err != nil {
		return APIKey{}, err
	}

	return key, nil
}

// ============================================================
// REVOKE
// ============================================================

func (r *apiKeyRepository) Revoke(
	ctx context.Context,
	userID string,
	id string,
) error {

	result, err := r.db.Exec(ctx, `
		UPDATE api_keys
		SET
			status = 'revoked',
			revoked_at = NOW()
		WHERE id = $1
		AND user_id = $2
		AND status = 'active'
	`,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("api key not found or already revoked")
	}

	return nil
}

// ============================================================
// GET BY PUBLISHABLE KEY
// ============================================================

func (r *apiKeyRepository) GetByPublishableKey(
	ctx context.Context,
	key string,
) (APIKey, error) {

	var apiKey APIKey

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status,
			created_at,
			revoked_at
		FROM api_keys
		WHERE publishable_key = $1
	`,
		key,
	).Scan(
		&apiKey.ID,
		&apiKey.UserID,
		&apiKey.Name,
		&apiKey.PublishableKey,
		&apiKey.SecretKeyHash,
		&apiKey.Status,
		&apiKey.CreatedAt,
		&apiKey.RevokedAt,
	)

	return apiKey, err
}

// ============================================================
// GET BY SECRET HASH
// ============================================================

func (r *apiKeyRepository) GetBySecretHash(
	ctx context.Context,
	hash string,
) (*APIKey, error) {

	var key APIKey

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			name,
			publishable_key,
			secret_key_hash,
			status,
			created_at,
			revoked_at
		FROM api_keys
		WHERE secret_key_hash = $1
		LIMIT 1
	`,
		hash,
	).Scan(
		&key.ID,
		&key.UserID,
		&key.Name,
		&key.PublishableKey,
		&key.SecretKeyHash,
		&key.Status,
		&key.CreatedAt,
		&key.RevokedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &key, nil
}
