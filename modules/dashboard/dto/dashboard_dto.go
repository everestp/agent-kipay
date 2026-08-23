// modules/dashboard/dto/dashboard_dto.go

package dto

type DashboardResponse struct {
	TotalBalance       float64 `json:"total_balance"`
	AvailableBalance   float64 `json:"available_balance"`
	ReservedBalance    float64 `json:"reserved_balance"`
	AgentSpendingToday float64 `json:"agent_spending_today"`
	Transactions       int64   `json:"transactions"`
	ActiveAgents       int64   `json:"active_agents"`
	APIPayments        int64   `json:"api_payments"`
	SuccessfulPayments int64   `json:"successful_payments"`
	FailedPayments     int64   `json:"failed_payments"`
	PendingPayments    int64   `json:"pending_payments"`
}
