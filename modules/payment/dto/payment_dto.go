package dto

type CreatePaymentRequest struct {
	AgentID        string  `json:"agent_id"`
	SessionID      string  `json:"session_id"`
	ServiceID      string  `json:"service_id"`
	Amount         float64 `json:"amount"`
	Asset          string  `json:"asset"`
	Network        string  `json:"network"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type PaymentResponse struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	Amount         float64 `json:"amount"`
	Asset          string  `json:"asset"`
	Network        string  `json:"network"`
	Protocol       string  `json:"protocol"`
	PolicyDecision string  `json:"policy_decision"`
	PolicyReason   string  `json:"policy_reason"`
	TxHash         *string `json:"tx_hash,omitempty"`
}
