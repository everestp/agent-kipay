// modules/user/controller/user_controller.go

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/user/dto"
	"github.com/everest/bheri/modules/user/service"
	"github.com/everest/bheri/pkg/utils"
)

type UserController struct {
	service service.UserService
}

func NewUserController(
	service service.UserService,
) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	user, err := c.service.Register(
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
		user,
	)
}

func (c *UserController) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	user, err := c.service.Login(
		r.Context(),
		request,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusUnauthorized,
			"invalid credentials",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		dto.LoginResponse{
			User: dto.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Status:    user.Status,
				CreatedAt: user.CreatedAt.String(),
			},
		},
	)
}

func (c *UserController) Me(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id")

	if userID == nil {
		utils.Error(
			w,
			http.StatusUnauthorized,
			"unauthorized",
		)
		return
	}

	user, err := c.service.GetByID(
		r.Context(),
		userID.(string),
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusNotFound,
			"user not found",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		user,
	)
}
