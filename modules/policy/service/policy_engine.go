package service

import (
	"context"
	"errors"
	"strings"

	"github.com/everest/bheri/modules/policy/dto"
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
}

func NewPolicyEngine(
	policyService PolicyService,
) PolicyEngine {
	return &policyEngine{
		policyService: policyService,
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

	if agentID == "" {
		return errors.New("agent id is required")
	}

	if sessionID == "" {
		return errors.New("session id is required")
	}

	if amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	if strings.TrimSpace(asset) == "" {
		return errors.New("asset is required")
	}

	if strings.TrimSpace(network) == "" {
		return errors.New("network is required")
	}

	// =========================================================
	// POLICY EVALUATION
	// =========================================================

	result, err := e.policyService.Evaluate(
		ctx,
		agentID,
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

	if result == nil {
		return errors.New(
			"policy evaluation returned no result",
		)
	}

	// =========================================================
	// POLICY DECISION
	// =========================================================

	if strings.ToLower(result.Decision) != "allowed" {

		if strings.TrimSpace(result.Reason) != "" {
			return errors.New(result.Reason)
		}

		return errors.New(
			"payment blocked by policy",
		)
	}

	return nil
}
