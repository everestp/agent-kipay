package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/session/dto"
	"github.com/everest/bheri/modules/session/repository"
)

type SessionService interface {
	Create(
		ctx context.Context,
		agentID string,
		request dto.CreateSessionRequest,
	) (*dto.SessionResponse, error)

	List(
		ctx context.Context,
		agentID string,
	) ([]models.Session, error)

	Revoke(
		ctx context.Context,
		id string,
	) error

	Validate(
		ctx context.Context,
		key string,
	) (*models.Session, error)
}

type sessionService struct {
	repository repository.SessionRepository
}

func NewSessionService(
	repository repository.SessionRepository,
) SessionService {
	return &sessionService{
		repository: repository,
	}
}

func generateSessionKey() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "bhs_" + hex.EncodeToString(bytes), nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))

	return hex.EncodeToString(hash[:])
}

func (s *sessionService) Create(
	ctx context.Context,
	agentID string,
	request dto.CreateSessionRequest,
) (*dto.SessionResponse, error) {

	if agentID == "" {
		return nil, errors.New("agent id is required")
	}

	if request.Name == "" {
		return nil, errors.New("session name is required")
	}

	if request.Limit <= 0 {
		return nil, errors.New(
			"session limit must be greater than zero",
		)
	}

	if request.ExpirationDays <= 0 {
		request.ExpirationDays = 1
	}

	key, err := generateSessionKey()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(
		time.Duration(request.ExpirationDays) *
			24 *
			time.Hour,
	)

	session, err := s.repository.Create(
		ctx,
		agentID,
		request.Name,
		hashKey(key),
		request.Limit,
		expiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &dto.SessionResponse{
		ID:        session.ID,
		Name:      session.Name,
		Status:    session.Status,
		Limit:     session.Limit,
		Spent:     session.Spent,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		Key:       key,
		Assets:    request.Assets,
		Networks:  request.Networks,
		Services:  request.Services,
	}, nil
}

func (s *sessionService) List(
	ctx context.Context,
	agentID string,
) ([]models.Session, error) {

	if agentID == "" {
		return nil, errors.New("agent id is required")
	}

	return s.repository.FindByAgentID(
		ctx,
		agentID,
	)
}

func (s *sessionService) Revoke(
	ctx context.Context,
	id string,
) error {

	if id == "" {
		return errors.New("session id is required")
	}

	return s.repository.Revoke(
		ctx,
		id,
	)
}

func (s *sessionService) Validate(
	ctx context.Context,
	key string,
) (*models.Session, error) {

	if key == "" {
		return nil, errors.New("session key is required")
	}

	session, err := s.repository.FindByKeyHash(
		ctx,
		hashKey(key),
	)

	if err != nil {
		return nil, errors.New("invalid session key")
	}

	if session.Status != "active" {
		return nil, errors.New("session is not active")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repository.Revoke(
			ctx,
			session.ID,
		)

		return nil, errors.New("session has expired")
	}

	if session.Spent >= session.Limit {
		return nil, errors.New(
			"session spending limit reached",
		)
	}

	return session, nil
}
