// modules/api_service/controller/api_service_controller.go

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/api_service/dto"
	"github.com/everest/bheri/modules/api_service/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type APIServiceController struct {
	service service.APIServiceService
}

func NewAPIServiceController(
	service service.APIServiceService,
) *APIServiceController {
	return &APIServiceController{
		service: service,
	}
}

func (c *APIServiceController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.CreateAPIServiceRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	service, err := c.service.Create(
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
		service,
	)
}

func (c *APIServiceController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	service, err := c.service.Get(
		r.Context(),
		chi.URLParam(r, "id"),
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusNotFound,
			"API service not found",
		)

		return
	}

	utils.Success(
		w,
		http.StatusOK,
		service,
	)
}

func (c *APIServiceController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	services, err := c.service.List(
		r.Context(),
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to fetch API services",
		)

		return
	}

	utils.Success(
		w,
		http.StatusOK,
		services,
	)
}

func (c *APIServiceController) ListActive(
	w http.ResponseWriter,
	r *http.Request,
) {

	services, err := c.service.ListActive(
		r.Context(),
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to fetch active API services",
		)

		return
	}

	utils.Success(
		w,
		http.StatusOK,
		services,
	)
}

func (c *APIServiceController) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.UpdateAPIServiceRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {

		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)

		return
	}

	service, err := c.service.Update(
		r.Context(),
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
		service,
	)
}

func (c *APIServiceController) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	err := c.service.Delete(
		r.Context(),
		chi.URLParam(r, "id"),
	)

	if err != nil {

		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to delete API service",
		)

		return
	}

	utils.Success(
		w,
		http.StatusOK,
		map[string]string{
			"message": "API service deleted",
		},
	)
}
