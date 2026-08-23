package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/settlement/dto"
	"github.com/everest/bheri/modules/settlement/service"

	"github.com/go-chi/chi/v5"
)

type SettlementController struct {
	service service.SettlementService
}

func NewSettlementController(
	service service.SettlementService,
) *SettlementController {
	return &SettlementController{
		service: service,
	}
}

func (c *SettlementController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	var req dto.CreateSettlementRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	result, err := c.service.Create(
		r.Context(),
		userID.(string),
		req,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(result)
}

func (c *SettlementController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	id := chi.URLParam(
		r,
		"id",
	)

	result, err := c.service.Get(
		r.Context(),
		userID.(string),
		id,
	)

	if err != nil {
		http.Error(
			w,
			"settlement not found",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(result)
}

func (c *SettlementController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	result, err := c.service.List(
		r.Context(),
		userID.(string),
	)

	if err != nil {
		http.Error(
			w,
			"failed to list settlements",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(result)
}

func (c *SettlementController) Verify(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	id := chi.URLParam(
		r,
		"id",
	)

	result, err := c.service.Verify(
		r.Context(),
		userID.(string),
		id,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(result)
}
