package dto

import "time"

// ============================================================
// CREATE API KEY
// ============================================================

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PublishableKey string    `json:"publishable_key"`
	SecretKey      string    `json:"secret_key"`
	CreatedAt      time.Time `json:"created_at"`
}

// ============================================================
// API KEY RESPONSE
// ============================================================

type APIKeyResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	PublishableKey string    `json:"publishable_key"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
}

// ============================================================
// REVOKE
// ============================================================

type RevokeAPIKeyRequest struct {
	Reason string `json:"reason"`
}
