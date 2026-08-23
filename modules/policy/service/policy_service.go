package service

import (
	"context"
	"errors"
	"time"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/policy/dto"
	"github.com/everest/bheri/modules/policy/repository"
)

type PolicyService interface {
	Create(
		ctx context.Context,
		request dto.CreatePolicyRequest,
	) (*models.Policy, error)

	Get(
		ctx context.Context,
		agentID string,
	) (*models.Policy, error)

	Evaluate(
		ctx context.Context,
		agentID string,
		request dto.EvaluatePolicyRequest,
	) (*dto.PolicyDecisionResponse, error)
}

type policyService struct {
	repository repository.PolicyRepository
}

func NewPolicyService(
	repository repository.PolicyRepository,
) PolicyService {
	return &policyService{
		repository: repository,
	}
}

func (s *policyService) Create(
	ctx context.Context,
	request dto.CreatePolicyRequest,
) (*models.Policy, error) {

	if request.AgentID == "" {
		return nil, errors.New("agent id is required")
	}

	if request.PerTransactionLimit <= 0 {
		return nil, errors.New(
			"per transaction limit must be greater than zero",
		)
	}

	if request.DailyLimit <= 0 {
		return nil, errors.New(
			"daily limit must be greater than zero",
		)
	}

	if request.WeeklyLimit <= 0 {
		return nil, errors.New(
			"weekly limit must be greater than zero",
		)
	}

	if request.ExpirationDays <= 0 {
		request.ExpirationDays = 30
	}

	return s.repository.Create(
		ctx,
		request.AgentID,
		repository.PolicyCreateData{
			DailyLimit:                  request.DailyLimit,
			PerTransactionLimit:         request.PerTransactionLimit,
			WeeklyLimit:                 request.WeeklyLimit,
			RequireApprovalAbove:        request.RequireApprovalAbove,
			ExpirationDays:              request.ExpirationDays,
			RequireApprovalNewMerchants: request.RequireApprovalNewMerchants,
			RequireApprovalNewAssets:    request.RequireApprovalNewAssets,
			BlockUnknownAPIs:            request.BlockUnknownAPIs,
			AutoPayments:                request.AutoPayments,
		},
	)
}

func (s *policyService) Get(
	ctx context.Context,
	agentID string,
) (*models.Policy, error) {

	return s.repository.FindByAgentID(
		ctx,
		agentID,
	)
}

func (s *policyService) Evaluate(
	ctx context.Context,
	agentID string,
	request dto.EvaluatePolicyRequest,
) (*dto.PolicyDecisionResponse, error) {

	policy, err := s.repository.FindByAgentID(
		ctx,
		agentID,
	)

	if err != nil {
		return nil, errors.New("active policy not found")
	}

	if policy.ExpiresAt != nil &&
		time.Now().After(*policy.ExpiresAt) {

		_ = s.repository.UpdateStatus(
			ctx,
			policy.ID,
			"expired",
		)

		return &dto.PolicyDecisionResponse{
			Decision: "blocked",
			Reason:   "policy has expired",
		}, nil
	}

	if request.Amount <= 0 {
		return &dto.PolicyDecisionResponse{
			Decision: "blocked",
			Reason:   "invalid payment amount",
		}, nil
	}

	if request.Amount > policy.PerTransactionLimit {
		return &dto.PolicyDecisionResponse{
			Decision: "blocked",
			Reason:   "per transaction limit exceeded",
		}, nil
	}

	if policy.RequireApprovalAbove > 0 &&
		request.Amount >= policy.RequireApprovalAbove {

		return &dto.PolicyDecisionResponse{
			Decision: "approval_required",
			Reason:   "payment requires approval",
		}, nil
	}

	if policy.BlockUnknownAPIs &&
		request.Service == "" {

		return &dto.PolicyDecisionResponse{
			Decision: "blocked",
			Reason:   "unknown API service",
		}, nil
	}

	return &dto.PolicyDecisionResponse{
		Decision: "allowed",
		Reason:   "payment satisfies policy",
	}, nil
}
