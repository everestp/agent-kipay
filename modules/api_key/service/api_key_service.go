package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/everest/bheri/modules/api_key/dto"
	"github.com/everest/bheri/modules/api_key/repository"
)

type APIKeyService interface {
	Create(
		ctx context.Context,
		userID string,
		req dto.CreateAPIKeyRequest,
	) (*dto.CreateAPIKeyResponse, error)

	List(
		ctx context.Context,
		userID string,
	) ([]dto.APIKeyResponse, error)

	Get(
		ctx context.Context,
		userID string,
		id string,
	) (*dto.APIKeyResponse, error)

	Revoke(
		ctx context.Context,
		userID string,
		id string,
	) error

	GetByPublishableKey(
		ctx context.Context,
		key string,
	) (repository.APIKey, error)

	Validate(
		ctx context.Context,
		secretKey string,
	) (*repository.APIKey, error)
}

type apiKeyService struct {
	repository repository.APIKeyRepository
}

func NewAPIKeyService(
	repository repository.APIKeyRepository,
) APIKeyService {
	return &apiKeyService{
		repository: repository,
	}
}

// ============================================================
// CREATE
// ============================================================

func (s *apiKeyService) Create(
	ctx context.Context,
	userID string,
	req dto.CreateAPIKeyRequest,
) (*dto.CreateAPIKeyResponse, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("api key name is required")
	}

	publishableKey, err := generateKey("bheri_pk_")
	if err != nil {
		return nil, err
	}

	secretKey, err := generateKey("bheri_sk_")
	if err != nil {
		return nil, err
	}

	secretKeyHash := hashSecret(secretKey)

	key, err := s.repository.Create(
		ctx,
		userID,
		req.Name,
		publishableKey,
		secretKeyHash,
	)

	if err != nil {
		return nil, err
	}

	return &dto.CreateAPIKeyResponse{
		ID:             key.ID,
		Name:           key.Name,
		PublishableKey: key.PublishableKey,
		SecretKey:      secretKey,
		CreatedAt:      key.CreatedAt,
	}, nil
}

// ============================================================
// LIST
// ============================================================

func (s *apiKeyService) List(
	ctx context.Context,
	userID string,
) ([]dto.APIKeyResponse, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	keys, err := s.repository.List(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	result := make(
		[]dto.APIKeyResponse,
		0,
		len(keys),
	)

	for _, key := range keys {
		result = append(
			result,
			dto.APIKeyResponse{
				ID:             key.ID,
				Name:           key.Name,
				PublishableKey: key.PublishableKey,
				Status:         key.Status,
				CreatedAt:      key.CreatedAt,
				RevokedAt:      key.RevokedAt,
			},
		)
	}

	return result, nil
}

// ============================================================
// GET
// ============================================================

func (s *apiKeyService) Get(
	ctx context.Context,
	userID string,
	id string,
) (*dto.APIKeyResponse, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("api key id is required")
	}

	key, err := s.repository.GetByID(
		ctx,
		userID,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &dto.APIKeyResponse{
		ID:             key.ID,
		Name:           key.Name,
		PublishableKey: key.PublishableKey,
		Status:         key.Status,
		CreatedAt:      key.CreatedAt,
		RevokedAt:      key.RevokedAt,
	}, nil
}

// ============================================================
// REVOKE
// ============================================================

func (s *apiKeyService) Revoke(
	ctx context.Context,
	userID string,
	id string,
) error {

	if strings.TrimSpace(userID) == "" {
		return errors.New("user id is required")
	}

	if strings.TrimSpace(id) == "" {
		return errors.New("api key id is required")
	}

	return s.repository.Revoke(
		ctx,
		userID,
		id,
	)
}

// ============================================================
// GET BY PUBLISHABLE KEY
// ============================================================

func (s *apiKeyService) GetByPublishableKey(
	ctx context.Context,
	key string,
) (repository.APIKey, error) {

	key = strings.TrimSpace(key)

	if key == "" {
		return repository.APIKey{},
			errors.New("publishable key is required")
	}

	if !strings.HasPrefix(key, "bheri_pk_") {
		return repository.APIKey{},
			errors.New("invalid publishable key format")
	}

	apiKey, err := s.repository.GetByPublishableKey(
		ctx,
		key,
	)

	if err != nil {
		return repository.APIKey{}, err
	}

	if apiKey.Status != "active" {
		return repository.APIKey{},
			errors.New("API key is revoked")
	}

	return apiKey, nil
}

// ============================================================
// VALIDATE SECRET KEY
// ============================================================

func (s *apiKeyService) Validate(
	ctx context.Context,
	secretKey string,
) (*repository.APIKey, error) {

	secretKey = strings.TrimSpace(secretKey)

	if secretKey == "" {
		return nil, errors.New("api key is required")
	}

	if !strings.HasPrefix(secretKey, "bheri_sk_") {
		return nil, errors.New("invalid API key format")
	}

	hash := hashSecret(secretKey)

	key, err := s.repository.GetBySecretHash(
		ctx,
		hash,
	)

	if err != nil {
		return nil, errors.New("invalid API key")
	}

	if key == nil {
		return nil, errors.New("invalid API key")
	}

	if key.Status != "active" {
		return nil, errors.New("API key is revoked")
	}

	return key, nil
}

// ============================================================
// GENERATE KEY
// ============================================================

func generateKey(prefix string) (string, error) {

	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return prefix + hex.EncodeToString(bytes), nil
}

// ============================================================
// HASH SECRET
// ============================================================

func hashSecret(secret string) string {

	hash := sha256.Sum256(
		[]byte(secret),
	)

	return hex.EncodeToString(hash[:])
}
