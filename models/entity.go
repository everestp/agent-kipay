// models/entity.go

package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Network   string    `json:"network"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Agent struct {
	ID                       string    `json:"id"`
	UserID                   string    `json:"user_id"`
	WalletID                 string    `json:"wallet_id"`
	Name                     string    `json:"name"`
	Description              string    `json:"description"`
	Status                   string    `json:"status"`
	Asset                    string    `json:"asset"`
	Network                  string    `json:"network"`
	Color                    string    `json:"color"`
	AutoPayments             bool      `json:"auto_payments"`
	LastActiveAt             *time.Time `json:"last_active_at"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type Policy struct {
	ID                          string    `json:"id"`
	AgentID                     string    `json:"agent_id"`
	DailyLimit                  float64   `json:"daily_limit"`
	PerTransactionLimit        float64   `json:"per_transaction_limit"`
	WeeklyLimit                 float64   `json:"weekly_limit"`
	RequireApprovalAbove       float64   `json:"require_approval_above"`
	ExpirationDays              int       `json:"expiration_days"`
	RequireApprovalNewMerchants bool      `json:"require_approval_new_merchants"`
	RequireApprovalNewAssets    bool      `json:"require_approval_new_assets"`
	BlockUnknownAPIs            bool      `json:"block_unknown_apis"`
	AutoPayments                bool      `json:"auto_payments"`
	Status                      string    `json:"status"`
	ExpiresAt                   *time.Time `json:"expires_at"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}


type Session struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Limit     float64   `json:"limit"`
	Spent     float64   `json:"spent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
type Payment struct {
	ID             string     `json:"id"`
	AgentID        string     `json:"agent_id"`
	SessionID      *string    `json:"session_id"`
	ServiceID      *string    `json:"service_id"`
	Amount         float64    `json:"amount"`
	Asset          string     `json:"asset"`
	Network        string     `json:"network"`
	Protocol       string     `json:"protocol"`
	Status         string     `json:"status"`
	PolicyDecision string     `json:"policy_decision"`
	PolicyReason   string     `json:"policy_reason"`
	IdempotencyKey string     `json:"idempotency_key"`
	PaymentNonce   string     `json:"payment_nonce"`
	TxHash         *string    `json:"tx_hash"`
	CreatedAt      time.Time  `json:"created_at"`
	ApprovedAt     *time.Time `json:"approved_at"`
	SubmittedAt    *time.Time `json:"submitted_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	FailedAt       *time.Time `json:"failed_at"`
}
type APIService struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`
	Endpoint          string    `json:"endpoint"`
	Status          string    `json:"status"`
	PricePerRequest   float64   `json:"price_per_request"`
	Asset             string    `json:"asset"`
	Network           string    `json:"network"`
	Description       string    `json:"description"`
	ProviderReputation float64  `json:"provider_reputation"`
	Active            bool      `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
type LedgerAccount struct {
	ID        string    `json:"id"`
	AgentID   *string   `json:"agent_id,omitempty"`
	WalletID  *string   `json:"wallet_id,omitempty"`
	Asset     string    `json:"asset"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LedgerEntry struct {
	ID           string    `json:"id"`
	PaymentID    *string   `json:"payment_id,omitempty"`
	SettlementID *string   `json:"settlement_id,omitempty"`
	AccountID    string    `json:"account_id"`
	EntryType    string    `json:"entry_type"`
	Debit        float64   `json:"debit"`
	Credit       float64   `json:"credit"`
	BalanceAfter float64   `json:"balance_after"`
	Reference    string    `json:"reference"`
	CreatedAt    time.Time `json:"created_at"`
}

type Transaction struct {
	ID           string     `json:"id"`
	PaymentID    *string    `json:"payment_id,omitempty"`
	SettlementID *string    `json:"settlement_id,omitempty"`
	AgentID      string     `json:"agent_id"`
	ServiceID    *string    `json:"service_id,omitempty"`
	Amount       float64    `json:"amount"`
	Asset        string     `json:"asset"`
	Network      string     `json:"network"`
	Protocol     string     `json:"protocol"`
	Status       string     `json:"status"`
	TxHash       *string    `json:"tx_hash,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
