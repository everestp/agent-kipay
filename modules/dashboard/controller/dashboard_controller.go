// modules/dashboard/controller/dashboard_controller.go

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/dashboard/service"
)

type DashboardController struct {
	service service.DashboardService
}

func NewDashboardController(
	service service.DashboardService,
) *DashboardController {
	return &DashboardController{
		service: service,
	}
}

func (c *DashboardController) Get(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	stats, err := c.service.GetDashboard(
		userID.(string),
	)

	if err != nil {
		http.Error(
			w,
			"failed to get dashboard",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(stats)
}
