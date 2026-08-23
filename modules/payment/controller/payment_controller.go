package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/payment/dto"
	"github.com/everest/bheri/modules/payment/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type PaymentController struct {
	service service.PaymentService
}

func NewPaymentController(
	service service.PaymentService,
) *PaymentController {
	return &PaymentController{
		service: service,
	}
}

func (c *PaymentController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	payment, err := c.service.Create(
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
		payment,
	)
}

func (c *PaymentController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	payment, err := c.service.Get(
		r.Context(),
		chi.URLParam(r, "id"),
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusNotFound,
			"payment not found",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		payment,
	)
}
