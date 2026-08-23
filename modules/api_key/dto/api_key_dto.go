package dto

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	PublishableKey string `json:"publishable_key"`
	SecretKey      string `json:"secret_key"`
	CreatedAt     string `json:"created_at"`
}

type APIKeyResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PublishableKey string `json:"publishable_key"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}
