package dto

type CreatePolicyRequest struct {
	AgentID                     string   `json:"agent_id"`
	DailyLimit                  float64  `json:"daily_limit"`
	PerTransactionLimit        float64  `json:"per_transaction_limit"`
	WeeklyLimit                 float64  `json:"weekly_limit"`
	RequireApprovalAbove        float64  `json:"require_approval_above"`
	AllowedAssets               []string `json:"allowed_assets"`
	AllowedNetworks             []string `json:"allowed_networks"`
	AllowedMerchants            []string `json:"allowed_merchants"`
	BlockedServices             []string `json:"blocked_services"`
	ExpirationDays              int      `json:"expiration_days"`
	RequireApprovalNewMerchants bool     `json:"require_approval_new_merchants"`
	RequireApprovalNewAssets    bool     `json:"require_approval_new_assets"`
	BlockUnknownAPIs            bool     `json:"block_unknown_apis"`
	AutoPayments                bool     `json:"auto_payments"`
}

type EvaluatePolicyRequest struct {
	Amount         float64 `json:"amount"`
	Asset          string  `json:"asset"`
	Network        string  `json:"network"`
	Merchant       string  `json:"merchant"`
	Service        string  `json:"service"`
}

type PolicyDecisionResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}
