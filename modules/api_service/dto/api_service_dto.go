package dto

type CreateAPIServiceRequest struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Endpoint        string  `json:"endpoint"`
	PricePerRequest float64 `json:"price_per_request"`
	Asset           string  `json:"asset"`
	Network         string  `json:"network"`
	Description     string  `json:"description"`
}

type APIServiceResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Endpoint        string  `json:"endpoint"`
	PricePerRequest float64 `json:"price_per_request"`
	Asset           string  `json:"asset"`
	Network         string  `json:"network"`
	Description     string  `json:"description"`
	Active          bool    `json:"active"`
}
