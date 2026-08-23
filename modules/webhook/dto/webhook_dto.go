package dto

type WebhookRequest struct {
	EventID     string                 `json:"event_id"`
	EventType   string                 `json:"event_type"`
	Network     string                 `json:"network"`
	Signature   string                 `json:"signature"`
	Timestamp   int64                  `json:"timestamp"`
	Data        map[string]interface{} `json:"data"`
}

type WebhookResponse struct {
	ID        string `json:"id"`
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
