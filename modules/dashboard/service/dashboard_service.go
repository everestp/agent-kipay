// modules/dashboard/service/dashboard_service.go

package service

import (
	"github.com/everest/bheri/modules/dashboard/dto"
	"github.com/everest/bheri/modules/dashboard/repository"
)

type DashboardService interface {
	GetDashboard(userID string) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	repository repository.DashboardRepository
}

func NewDashboardService(
	repository repository.DashboardRepository,
) DashboardService {
	return &dashboardService{
		repository: repository,
	}
}

func (s *dashboardService) GetDashboard(
	userID string,
) (*dto.DashboardResponse, error) {

	stats, err := s.repository.GetStats(userID)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardResponse{
		TotalBalance:       stats.TotalBalance,
		AvailableBalance:   stats.AvailableBalance,
		ReservedBalance:    stats.ReservedBalance,
		AgentSpendingToday: stats.AgentSpendingToday,
		Transactions:       stats.Transactions,
		ActiveAgents:       stats.ActiveAgents,
		APIPayments:        stats.APIPayments,
		SuccessfulPayments: stats.SuccessfulPayments,
		FailedPayments:     stats.FailedPayments,
		PendingPayments:    stats.PendingPayments,
	}, nil
}
