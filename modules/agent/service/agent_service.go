package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/agent/dto"
	"github.com/everest/bheri/modules/agent/repository"
)

type AgentService interface {
	Create(
		ctx context.Context,
		userID string,
		request dto.CreateAgentRequest,
	) (*models.Agent, error)

	Get(
		ctx context.Context,
		userID string,
		id string,
	) (*models.Agent, error)

	List(
		ctx context.Context,
		userID string,
	) ([]models.Agent, error)

	Update(
		ctx context.Context,
		userID string,
		id string,
		request dto.UpdateAgentRequest,
	) (*models.Agent, error)

	Delete(
		ctx context.Context,
		userID string,
		id string,
	) error
}

type agentService struct {
	repository repository.AgentRepository
}

func NewAgentService(
	repository repository.AgentRepository,
) AgentService {
	return &agentService{
		repository: repository,
	}
}

func (s *agentService) Create(
	ctx context.Context,
	userID string,
	request dto.CreateAgentRequest,
) (*models.Agent, error) {

	if request.Name == "" {
		return nil, errors.New("agent name is required")
	}

	if request.WalletID == "" {
		return nil, errors.New("wallet id is required")
	}

	if request.Asset == "" {
		request.Asset = "USDC"
	}

	if request.Network == "" {
		request.Network = "solana-devnet"
	}

	return s.repository.Create(
		ctx,
		userID,
		request.WalletID,
		request.Name,
		request.Description,
		request.Asset,
		request.Network,
		request.Color,
	)
}

func (s *agentService) Get(
	ctx context.Context,
	userID string,
	id string,
) (*models.Agent, error) {

	return s.repository.FindByID(
		ctx,
		userID,
		id,
	)
}

func (s *agentService) List(
	ctx context.Context,
	userID string,
) ([]models.Agent, error) {

	return s.repository.FindByUserID(
		ctx,
		userID,
	)
}

func (s *agentService) Update(
	ctx context.Context,
	userID string,
	id string,
	request dto.UpdateAgentRequest,
) (*models.Agent, error) {

	if request.Name == "" {
		return nil, errors.New("agent name is required")
	}

	if request.Status == "" {
		request.Status = "idle"
	}

	autoPayments := true

	if request.AutoPayments != nil {
		autoPayments = *request.AutoPayments
	}

	return s.repository.Update(
		ctx,
		userID,
		id,
		request.Name,
		request.Description,
		request.Status,
		request.Color,
		autoPayments,
	)
}

func (s *agentService) Delete(
	ctx context.Context,
	userID string,
	id string,
) error {

	return s.repository.Delete(
		ctx,
		userID,
		id,
	)
}
