package dto

type CreateAgentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WalletID    string `json:"wallet_id"`
	Asset       string `json:"asset"`
	Network     string `json:"network"`
	Color       string `json:"color"`
}

type UpdateAgentRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	Color         string `json:"color"`
	AutoPayments  *bool  `json:"auto_payments"`
}
