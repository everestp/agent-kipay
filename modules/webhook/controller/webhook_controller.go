package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/webhook/dto"
	"github.com/everest/bheri/modules/webhook/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type WebhookController struct {
	service *service.WebhookService
}

func NewWebhookController(
	service *service.WebhookService,
) *WebhookController {

	return &WebhookController{
		service: service,
	}
}

func (c *WebhookController) Receive(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.WebhookRequest

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

	result, err := c.service.Process(
		r.Context(),
		req,
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusInternalServerError,
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

func (c *WebhookController) Get(
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
			"webhook not found",
		)

		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		result,
	)
}
