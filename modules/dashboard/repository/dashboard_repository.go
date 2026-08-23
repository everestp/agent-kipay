// modules/dashboard/repository/dashboard_repository.go

package repository

import (
	"database/sql"
)

type DashboardRepository interface {
	GetStats(userID string) (DashboardStats, error)
}

type DashboardStats struct {
	TotalBalance       float64
	AvailableBalance   float64
	ReservedBalance    float64
	AgentSpendingToday float64
	Transactions       int64
	ActiveAgents       int64
	APIPayments        int64
	SuccessfulPayments int64
	FailedPayments     int64
	PendingPayments    int64
}

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

func (r *dashboardRepository) GetStats(userID string) (DashboardStats, error) {

	var stats DashboardStats

	err := r.db.QueryRow(`
		SELECT
			COALESCE((
				SELECT SUM(w.balance)
				FROM wallets w
				WHERE w.user_id = $1
			), 0),

			COALESCE((
				SELECT SUM(w.available_balance)
				FROM wallets w
				WHERE w.user_id = $1
			), 0),

			COALESCE((
				SELECT SUM(w.reserved_balance)
				FROM wallets w
				WHERE w.user_id = $1
			), 0),

			COALESCE((
				SELECT SUM(p.amount)
				FROM payments p
				JOIN agents a ON a.id = p.agent_id
				WHERE a.user_id = $1
				AND p.created_at >= CURRENT_DATE
				AND p.status = 'completed'
			), 0),

			(
				SELECT COUNT(*)
				FROM transactions t
				JOIN agents a ON a.id = t.agent_id
				WHERE a.user_id = $1
			),

			(
				SELECT COUNT(*)
				FROM agents
				WHERE user_id = $1
				AND status = 'active'
			),

			(
				SELECT COUNT(*)
				FROM payments p
				JOIN agents a ON a.id = p.agent_id
				WHERE a.user_id = $1
				AND p.protocol = 'x402'
			),

			(
				SELECT COUNT(*)
				FROM payments p
				JOIN agents a ON a.id = p.agent_id
				WHERE a.user_id = $1
				AND p.status = 'completed'
			),

			(
				SELECT COUNT(*)
				FROM payments p
				JOIN agents a ON a.id = p.agent_id
				WHERE a.user_id = $1
				AND p.status = 'failed'
			),

			(
				SELECT COUNT(*)
				FROM payments p
				JOIN agents a ON a.id = p.agent_id
				WHERE a.user_id = $1
				AND p.status = 'pending'
			)
	`,
		userID,
	).Scan(
		&stats.TotalBalance,
		&stats.AvailableBalance,
		&stats.ReservedBalance,
		&stats.AgentSpendingToday,
		&stats.Transactions,
		&stats.ActiveAgents,
		&stats.APIPayments,
		&stats.SuccessfulPayments,
		&stats.FailedPayments,
		&stats.PendingPayments,
	)

	if err != nil {
		return DashboardStats{}, err
	}

	return stats, nil
}
