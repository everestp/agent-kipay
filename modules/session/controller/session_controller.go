package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/session/dto"
	"github.com/everest/bheri/modules/session/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type SessionController struct {
	service service.SessionService
}

func NewSessionController(
	service service.SessionService,
) *SessionController {
	return &SessionController{
		service: service,
	}
}

func (c *SessionController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	agentID := chi.URLParam(r, "agentID")

	var request dto.CreateSessionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	session, err := c.service.Create(
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
		http.StatusCreated,
		session,
	)
}

func (c *SessionController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	agentID := chi.URLParam(r, "agentID")

	sessions, err := c.service.List(
		r.Context(),
		agentID,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to fetch sessions",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		sessions,
	)
}

func (c *SessionController) Revoke(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	err := c.service.Revoke(
		r.Context(),
		id,
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
		map[string]string{
			"message": "session revoked",
		},
	)
}
