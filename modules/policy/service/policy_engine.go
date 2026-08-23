package service

import (
	"context"
	"errors"
)

type PolicyEngine interface {
	Authorize(
		ctx context.Context,
		agentID string,
		sessionID string,
		amount float64,
		asset string,
		network string,
		service string,
		merchant string,
	) error
}

type policyEngine struct {
	policyService PolicyService
	sessionService SessionService
}

func NewPolicyEngine(
	policyService PolicyService,
	sessionService SessionService,
) PolicyEngine {
	return &policyEngine{
		policyService:  policyService,
		sessionService: sessionService,
	}
}

func (e *policyEngine) Authorize(
	ctx context.Context,
	agentID string,
	sessionID string,
	amount float64,
	asset string,
	network string,
	service string,
	merchant string,
) error {

	if amount <= 0 {
		return errors.New(
			"payment amount must be greater than zero",
		)
	}

	session, err := e.sessionService.repository.
		FindByKeyHash(ctx, sessionID)

	_ = session
	_ = err

	result, err := e.policyService.Evaluate(
		ctx,
		agentID,
		// request created by caller
		dto.EvaluatePolicyRequest{
			Amount:   amount,
			Asset:    asset,
			Network:  network,
			Service:  service,
			Merchant: merchant,
		},
	)

	if err != nil {
		return err
	}

	if result.Decision != "allowed" {
		return errors.New(
			result.Reason,
		)
	}

	return nil
}
