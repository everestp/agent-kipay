package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/everest/bheri/modules/ledger/dto"
	"github.com/everest/bheri/modules/ledger/service"
	"github.com/go-chi/chi/v5"
)

type LedgerController struct {
	service *service.LedgerService
}

func NewLedgerController(
	service *service.LedgerService,
) *LedgerController {

	return &LedgerController{
		service: service,
	}
}

// ============================================================
// CREATE ACCOUNT
// ============================================================

func (c *LedgerController) CreateAccount(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateLedgerAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.Name == "" {
		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)
		return
	}

	if req.Currency == "" {
		http.Error(
			w,
			"currency is required",
			http.StatusBadRequest,
		)
		return
	}

	if req.WalletID == nil && req.AgentID == nil {
		http.Error(
			w,
			"walletId or agentId is required",
			http.StatusBadRequest,
		)
		return
	}

	account, err := c.service.CreateAccount(
		r.Context(),
		req,
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

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(account)
}

// ============================================================
// GET ACCOUNT
// ============================================================

func (c *LedgerController) GetAccount(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	account, err := c.service.GetAccount(
		r.Context(),
		id,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			http.Error(
				w,
				"ledger account not found",
				http.StatusNotFound,
			)
			return
		}

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

	_ = json.NewEncoder(w).Encode(account)
}

// ============================================================
// LIST ACCOUNTS
// ============================================================

func (c *LedgerController) ListAccounts(
	w http.ResponseWriter,
	r *http.Request,
) {

	accounts, err := c.service.ListAccounts(
		r.Context(),
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

	_ = json.NewEncoder(w).Encode(accounts)
}

// ============================================================
// CREATE LEDGER TRANSACTION
// ============================================================

func (c *LedgerController) CreateTransaction(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateLedgerTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.DebitAccountID == "" {
		http.Error(
			w,
			"debitAccountId is required",
			http.StatusBadRequest,
		)
		return
	}

	if req.CreditAccountID == "" {
		http.Error(
			w,
			"creditAccountId is required",
			http.StatusBadRequest,
		)
		return
	}

	if req.DebitAccountID == req.CreditAccountID {
		http.Error(
			w,
			"debit and credit accounts must be different",
			http.StatusBadRequest,
		)
		return
	}

	if req.Amount <= 0 {
		http.Error(
			w,
			"amount must be greater than zero",
			http.StatusBadRequest,
		)
		return
	}

	if req.Asset == "" {
		http.Error(
			w,
			"asset is required",
			http.StatusBadRequest,
		)
		return
	}

	result, err := c.service.CreateTransaction(
		r.Context(),
		req,
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

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(result)
}

// ============================================================
// GET LEDGER TRANSACTION
// ============================================================

func (c *LedgerController) GetTransaction(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	result, err := c.service.GetTransaction(
		r.Context(),
		id,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			http.Error(
				w,
				"ledger transaction not found",
				http.StatusNotFound,
			)
			return
		}

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
// ACCOUNT ENTRIES
// ============================================================

func (c *LedgerController) ListAccountEntries(
	w http.ResponseWriter,
	r *http.Request,
) {

	accountID := chi.URLParam(
		r,
		"id",
	)

	entries, err := c.service.ListAccountEntries(
		r.Context(),
		accountID,
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

	_ = json.NewEncoder(w).Encode(entries)
}
