package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/policy/dto"
	"github.com/everest/bheri/modules/policy/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type PolicyController struct {
	service service.PolicyService
}

func NewPolicyController(
	service service.PolicyService,
) *PolicyController {
	return &PolicyController{
		service: service,
	}
}

func (c *PolicyController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreatePolicyRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	policy, err := c.service.Create(
		r.Context(),
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
		policy,
	)
}

func (c *PolicyController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	agentID := chi.URLParam(r, "agentID")

	policy, err := c.service.Get(
		r.Context(),
		agentID,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusNotFound,
			"active policy not found",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		policy,
	)
}

func (c *PolicyController) Evaluate(
	w http.ResponseWriter,
	r *http.Request,
) {

	agentID := chi.URLParam(r, "agentID")

	var request dto.EvaluatePolicyRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	result, err := c.service.Evaluate(
		r.Context(),
		agentID,
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
		result,
	)
}
