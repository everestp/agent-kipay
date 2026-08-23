package dto

type TransferRequest struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      float64 `json:"amount"`
	Asset       string  `json:"asset"`
	Network     string  `json:"network"`
}

type TransferResponse struct {
	TxHash string `json:"tx_hash"`
}

type VerifyTransactionRequest struct {
	TxHash  string `json:"tx_hash"`
	Network string `json:"network"`
}

type VerifyTransactionResponse struct {
	Confirmed     bool   `json:"confirmed"`
	BlockNumber   uint64 `json:"block_number"`
	Confirmations int    `json:"confirmations"`
}
