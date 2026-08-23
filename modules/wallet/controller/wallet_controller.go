// modules/wallet/controller/wallet_controller.go

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/wallet/dto"
	"github.com/everest/bheri/modules/wallet/service"
	"github.com/everest/bheri/pkg/utils"

	"github.com/go-chi/chi/v5"
)

type WalletController struct {
	service service.WalletService
}

func NewWalletController(
	service service.WalletService,
) *WalletController {
	return &WalletController{
		service: service,
	}
}

func (c *WalletController) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	var request dto.CreateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	wallet, err := c.service.Create(
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
		wallet,
	)
}

func (c *WalletController) Get(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	id := chi.URLParam(r, "id")

	wallet, err := c.service.Get(
		r.Context(),
		id,
		userID,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusNotFound,
			"wallet not found",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		wallet,
	)
}

func (c *WalletController) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID := r.Context().Value("user_id").(string)

	wallets, err := c.service.List(
		r.Context(),
		userID,
	)

	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"failed to fetch wallets",
		)
		return
	}

	utils.Success(
		w,
		http.StatusOK,
		wallets,
	)
}
