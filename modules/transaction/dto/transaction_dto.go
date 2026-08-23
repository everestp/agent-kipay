package dto

type CreateTransactionRequest struct {
	PaymentID      string  `json:"payment_id"`
	AgentID        string  `json:"agent_id"`
	SessionID      *string `json:"session_id,omitempty"`
	APIServiceID   *string `json:"api_service_id,omitempty"`
	Type           string  `json:"type"`
	Protocol       string  `json:"protocol"`
	Network        string  `json:"network"`
	Asset          string  `json:"asset"`
	Amount         float64 `json:"amount"`
	SenderAddress  string  `json:"sender_address"`
	ReceiverAddress string `json:"receiver_address"`
}

type UpdateTransactionRequest struct {
	Status           string  `json:"status"`
	TxHash           *string `json:"tx_hash,omitempty"`
	BlockNumber      *int64  `json:"block_number,omitempty"`
	VerificationStatus string `json:"verification_status,omitempty"`
	SettlementStatus   string `json:"settlement_status,omitempty"`
}

type TransactionResponse struct {
	ID                 string  `json:"id"`
	PaymentID          string  `json:"payment_id"`
	AgentID            string  `json:"agent_id"`
	SessionID          *string `json:"session_id,omitempty"`
	APIServiceID       *string `json:"api_service_id,omitempty"`
	Type               string  `json:"type"`
	Protocol           string  `json:"protocol"`
	Network            string  `json:"network"`
	Asset              string  `json:"asset"`
	Amount             float64 `json:"amount"`
	SenderAddress      string  `json:"sender_address"`
	ReceiverAddress    string  `json:"receiver_address"`
	Status             string  `json:"status"`
	TxHash             *string `json:"tx_hash,omitempty"`
	BlockNumber        *int64  `json:"block_number,omitempty"`
	VerificationStatus string  `json:"verification_status"`
	SettlementStatus   string  `json:"settlement_status"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}
