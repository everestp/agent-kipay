package dto


type CreateLedgerAccountRequest struct {
	WalletID *string `json:"walletId"`
	AgentID  *string `json:"agentId"`
	Name     string  `json:"name"`
	Currency string `json:"currency"`
}

type LedgerAccountResponse struct {
	ID        string  `json:"id"`
	WalletID  *string `json:"walletId,omitempty"`
	AgentID   *string `json:"agentId,omitempty"`
	Name      string  `json:"name"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"createdAt"`
}

type LedgerEntryResponse struct {
	ID                 string  `json:"id"`
	LedgerTransactionID string  `json:"ledgerTransactionId"`
	LedgerAccountID    string  `json:"ledgerAccountId"`
	EntryType          string  `json:"entryType"`
	Amount             float64 `json:"amount"`
	Asset              string  `json:"asset"`
	CreatedAt          string  `json:"createdAt"`
}

type LedgerTransactionResponse struct {
	ID            string                 `json:"id"`
	PaymentID     *string                `json:"paymentId,omitempty"`
	TransactionID *string                `json:"transactionId,omitempty"`
	Reference     *string                `json:"reference,omitempty"`
	CreatedAt     string                 `json:"createdAt"`
	Entries       []LedgerEntryResponse  `json:"entries"`
}

type CreateLedgerTransactionRequest struct {
	PaymentID     *string `json:"paymentId"`
	TransactionID *string `json:"transactionId"`
	Reference     string  `json:"reference"`

	DebitAccountID  string  `json:"debitAccountId"`
	CreditAccountID string  `json:"creditAccountId"`
	Amount          float64 `json:"amount"`
	Asset           string  `json:"asset"`
}
