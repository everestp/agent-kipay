package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/api_key/dto"
	"github.com/everest/bheri/modules/api_key/service"

	"github.com/go-chi/chi/v5"
)

type APIKeyController struct {
	service service.APIKeyService
}

func NewAPIKeyController(
	service service.APIKeyService,
) *APIKeyController {
	return &APIKeyController{
		service: service,
	}
}

// ============================================================
// CREATE
// ============================================================

func (c *APIKeyController) Create(
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

	var req dto.CreateAPIKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	_ = json.NewEncoder(w).Encode(result)
}

// ============================================================
// LIST
// ============================================================

func (c *APIKeyController) List(
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
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(result)
}

// ============================================================
// GET
// ============================================================

func (c *APIKeyController) Get(
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
			err.Error(),
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	_ = json.NewEncoder(w).Encode(result)
}

// ============================================================
// REVOKE
// ============================================================

func (c *APIKeyController) Revoke(
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

	err := c.service.Revoke(
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

	_ = json.NewEncoder(w).Encode(
		map[string]interface{}{
			"success": true,
			"message": "api key revoked",
		},
	)
}
