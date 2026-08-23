package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/agent/dto"
	"github.com/everest/bheri/modules/agent/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type AgentController struct {
	service service.AgentService
}

func NewAgentController(
	service service.AgentService,
) *AgentController {
	return &AgentController{
		service: service,
	}
}

func (c *AgentController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	var request dto.CreateAgentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	agent, err := c.service.Create(
		r.Context(),
		userID,
		request,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	utils.Success(
		w,
		http.StatusCreated,
		agent,
	)
}

func (c *AgentController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	agent, err := c.service.Get(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusNotFound,
			"agent not found",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		agent,
	)
}

func (c *AgentController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	agents, err := c.service.List(
		r.Context(),
		userID,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to fetch agents",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		agents,
	)
}

func (c *AgentController) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	var request dto.UpdateAgentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	agent, err := c.service.Update(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
		request,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		agent,
	)
}

func (c *AgentController) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	err := c.service.Delete(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to delete agent",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		map[string]string{
			"message": "agent deleted",
		},
	)
}
