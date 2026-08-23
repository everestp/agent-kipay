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
