package dto

type CreateSettlementRequest struct {
	PaymentID         string  `json:"payment_id"`
	TxHash            string  `json:"tx_hash"`
	Network           string  `json:"network"`
	ExpectedAmount    float64 `json:"expected_amount"`
	Asset             string  `json:"asset"`
	ReceiverAddress   string  `json:"receiver_address"`
}

type SettlementResponse struct {
	ID            string  `json:"id"`
	PaymentID     string  `json:"payment_id"`
	TxHash        string  `json:"tx_hash"`
	Network       string  `json:"network"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Asset         string  `json:"asset"`
	Receiver      string  `json:"receiver"`
	BlockNumber   uint64  `json:"block_number"`
	Message       string  `json:"message"`
	CreatedAt     string  `json:"created_at"`
	ConfirmedAt   *string `json:"confirmed_at,omitempty"`
}
