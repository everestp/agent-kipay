
package dto

type CreateSettlementRequest struct {
	PaymentID string `json:"payment_id"`
}

type SettlementResponse struct {
	ID              string `json:"id"`
	PaymentID       string `json:"payment_id"`
	TxHash          *string `json:"tx_hash,omitempty"`
	Status          string `json:"status"`
	Confirmations   int    `json:"confirmations"`
	BlockNumber     *int64  `json:"block_number,omitempty"`
}
