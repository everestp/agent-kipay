package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/transaction/dto"
	"github.com/everest/bheri/modules/transaction/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type TransactionController struct {
	service *service.TransactionService
}

func NewTransactionController(
	service *service.TransactionService,
) *TransactionController {

	return &TransactionController{
		service: service,
	}
}

func (c *TransactionController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateTransactionRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	result, err := c.service.Create(
		r.Context(),
		req,
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	utils.JSON(
		w,
		http.StatusCreated,
		result,
	)
}

func (c *TransactionController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(
		r,
		"id",
	)

	result, err := c.service.Get(
		r.Context(),
		id,
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusNotFound,
			"transaction not found",
		)

		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		result,
	)
}

func (c *TransactionController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	agentID := r.URL.Query().Get(
		"agent_id",
	)

	result, err := c.service.List(
		r.Context(),
		agentID,
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		result,
	)
}

func (c *TransactionController) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(
		r,
		"id",
	)

	var req dto.UpdateTransactionRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	if err := c.service.UpdateStatus(
		r.Context(),
		id,
		req,
	); err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"message": "transaction updated successfully",
		},
	)
}
