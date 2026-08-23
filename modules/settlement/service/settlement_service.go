package service

import (
	"context"
	"errors"
	"strings"

	blockchaindto "github.com/everest/bheri/modules/blockchain/dto"
	blockchainservice "github.com/everest/bheri/modules/blockchain/service"
	"github.com/everest/bheri/modules/settlement/dto"
	"github.com/everest/bheri/modules/settlement/repository"
)

type SettlementService interface {
	Create(
		ctx context.Context,
		userID string,
		req dto.CreateSettlementRequest,
	) (*dto.SettlementResponse, error)

	Get(
		ctx context.Context,
		userID string,
		id string,
	) (*dto.SettlementResponse, error)

	List(
		ctx context.Context,
		userID string,
	) ([]dto.SettlementResponse, error)

	Verify(
		ctx context.Context,
		userID string,
		id string,
	) (*dto.SettlementResponse, error)
}

type settlementService struct {
	repository       repository.SettlementRepository
	blockchainService blockchainservice.BlockchainService
}

func NewSettlementService(
	repository repository.SettlementRepository,
	blockchainService blockchainservice.BlockchainService,
) SettlementService {
	return &settlementService{
		repository:        repository,
		blockchainService: blockchainService,
	}
}

// ============================================================
// CREATE
// ============================================================

func (s *settlementService) Create(
	ctx context.Context,
	userID string,
	req dto.CreateSettlementRequest,
) (*dto.SettlementResponse, error) {

	if ctx == nil {
		return nil, errors.New("context is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(req.PaymentID) == "" {
		return nil, errors.New("payment_id is required")
	}

	if strings.TrimSpace(req.TxHash) == "" {
		return nil, errors.New("tx_hash is required")
	}

	if strings.TrimSpace(req.Network) == "" {
		return nil, errors.New("network is required")
	}

	if strings.TrimSpace(req.Asset) == "" {
		return nil, errors.New("asset is required")
	}

	if strings.TrimSpace(req.ReceiverAddress) == "" {
		return nil, errors.New(
			"receiver_address is required",
		)
	}

	settlement := repository.Settlement{
		UserID:    userID,
		PaymentID: strings.TrimSpace(req.PaymentID),
		TxHash:    strings.TrimSpace(req.TxHash),
		Network:   strings.TrimSpace(req.Network),
		Status:    "pending",
		Amount:    0,
		Asset:     strings.TrimSpace(req.Asset),
		Receiver:  strings.TrimSpace(req.ReceiverAddress),
		BlockNumber: 0,
		Message:   "awaiting blockchain verification",
	}

	created, err := s.repository.Create(
		ctx,
		settlement,
	)
	if err != nil {
		return nil, err
	}

	return toResponse(created), nil
}

// ============================================================
// GET
// ============================================================

func (s *settlementService) Get(
	ctx context.Context,
	userID string,
	id string,
) (*dto.SettlementResponse, error) {

	if ctx == nil {
		return nil, errors.New("context is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("settlement id is required")
	}

	settlement, err := s.repository.GetByID(
		ctx,
		userID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return toResponse(settlement), nil
}

// ============================================================
// LIST
// ============================================================

func (s *settlementService) List(
	ctx context.Context,
	userID string,
) ([]dto.SettlementResponse, error) {

	if ctx == nil {
		return nil, errors.New("context is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	settlements, err := s.repository.List(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	result := make(
		[]dto.SettlementResponse,
		0,
		len(settlements),
	)

	for _, settlement := range settlements {
		result = append(
			result,
			*toResponse(settlement),
		)
	}

	return result, nil
}

// ============================================================
// VERIFY
// ============================================================

func (s *settlementService) Verify(
	ctx context.Context,
	userID string,
	id string,
) (*dto.SettlementResponse, error) {

	if ctx == nil {
		return nil, errors.New("context is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user id is required")
	}

	if strings.TrimSpace(id) == "" {
		return nil, errors.New("settlement id is required")
	}

	// ----------------------------------------------------------
	// Load settlement
	// ----------------------------------------------------------

	settlement, err := s.repository.GetByID(
		ctx,
		userID,
		id,
	)
	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Already confirmed
	// ----------------------------------------------------------

	if settlement.Status == "confirmed" {
		return toResponse(settlement), nil
	}

	// ----------------------------------------------------------
	// Already failed
	// ----------------------------------------------------------

	if settlement.Status == "failed" {
		return nil, errors.New(
			"settlement has already failed",
		)
	}

	// ----------------------------------------------------------
	// Validate settlement
	// ----------------------------------------------------------

	if strings.TrimSpace(settlement.TxHash) == "" {
		return nil, errors.New(
			"settlement transaction hash is missing",
		)
	}

	if strings.TrimSpace(settlement.Network) == "" {
		return nil, errors.New(
			"settlement network is missing",
		)
	}

	if strings.TrimSpace(settlement.Asset) == "" {
		return nil, errors.New(
			"settlement asset is missing",
		)
	}

	if strings.TrimSpace(settlement.Receiver) == "" {
		return nil, errors.New(
			"settlement receiver is missing",
		)
	}

	// ----------------------------------------------------------
	// Mark settlement as verifying
	// ----------------------------------------------------------

	err = s.repository.UpdateStatus(
		ctx,
		settlement.ID,
		"verifying",
		settlement.Amount,
		settlement.BlockNumber,
		"verifying blockchain transaction",
	)
	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Blockchain verification
	// ----------------------------------------------------------

	result, err := s.blockchainService.VerifyTransaction(
		ctx,
		blockchaindto.VerifyTransactionRequest{
			TxHash:           settlement.TxHash,
			Network:          settlement.Network,
			ExpectedAmount:   settlement.Amount,
			ExpectedAsset:    settlement.Asset,
			ExpectedReceiver: settlement.Receiver,
		},
	)

	if err != nil {

		updateErr := s.repository.UpdateStatus(
			ctx,
			settlement.ID,
			"failed",
			settlement.Amount,
			settlement.BlockNumber,
			err.Error(),
		)

		if updateErr != nil {
			return nil, errors.Join(
				err,
				updateErr,
			)
		}

		return nil, err
	}

	// ----------------------------------------------------------
	// Determine final status
	// ----------------------------------------------------------

	status := "failed"

	if result.Valid && result.Confirmed {
		status = "confirmed"
	}

	// ----------------------------------------------------------
	// Determine message
	// ----------------------------------------------------------

	message := strings.TrimSpace(
		result.Message,
	)

	if message == "" {
		if status == "confirmed" {
			message = "transaction confirmed"
		} else {
			message = "transaction verification failed"
		}
	}

	// ----------------------------------------------------------
	// Persist blockchain result
	// ----------------------------------------------------------

	err = s.repository.UpdateStatus(
		ctx,
		settlement.ID,
		status,
		result.Amount,
		result.BlockNumber,
		message,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Fetch updated settlement
	// ----------------------------------------------------------

	updated, err := s.repository.GetByID(
		ctx,
		userID,
		id,
	)
	if err != nil {
		return nil, err
	}

	return toResponse(updated), nil
}

// ============================================================
// RESPONSE MAPPER
// ============================================================

func toResponse(
	s repository.Settlement,
) *dto.SettlementResponse {

	return &dto.SettlementResponse{
		ID:           s.ID,
		PaymentID:    s.PaymentID,
		TxHash:       s.TxHash,
		Network:      s.Network,
		Status:       s.Status,
		Amount:       s.Amount,
		Asset:        s.Asset,
		Receiver:     s.Receiver,
		BlockNumber:  s.BlockNumber,
		Message:      s.Message,
		CreatedAt:    s.CreatedAt,
		ConfirmedAt:  s.ConfirmedAt,
	}
}
