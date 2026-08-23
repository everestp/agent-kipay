package dto

type CreateSessionRequest struct {
	Name           string   `json:"name"`
	Limit          float64  `json:"limit"`
	ExpirationDays int      `json:"expiration_days"`
	Assets         []string `json:"assets"`
	Networks       []string `json:"networks"`
	Services       []string `json:"services"`
}

type SessionResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Status    string   `json:"status"`
	Limit     float64  `json:"limit"`
	Spent     float64  `json:"spent"`
	ExpiresAt string   `json:"expires_at"`
	Key       string   `json:"key"`
	Assets    []string `json:"assets"`
	Networks  []string `json:"networks"`
	Services  []string `json:"services"`
}

type RevokeSessionRequest struct {
	Reason string `json:"reason"`
}
