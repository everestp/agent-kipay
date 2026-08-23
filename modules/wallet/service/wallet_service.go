// modules/wallet/service/wallet_service.go

package service

import (
	"context"
	"errors"

	"github.com/everest/bheri/models"
	"github.com/everest/bheri/modules/wallet/dto"
	"github.com/everest/bheri/modules/wallet/repository"
)

type WalletService interface {
	Create(
		ctx context.Context,
		userID string,
		request dto.CreateWalletRequest,
	) (*models.Wallet, error)

	Get(
		ctx context.Context,
		id string,
		userID string,
	) (*models.Wallet, error)

	List(
		ctx context.Context,
		userID string,
	) ([]models.Wallet, error)
}

type walletService struct {
	repository repository.WalletRepository
}

func NewWalletService(
	repository repository.WalletRepository,
) WalletService {
	return &walletService{
		repository: repository,
	}
}

func (s *walletService) Create(
	ctx context.Context,
	userID string,
	request dto.CreateWalletRequest,
) (*models.Wallet, error) {

	if request.Name == "" {
		request.Name = "Main Wallet"
	}

	if request.Address == "" {
		return nil, errors.New("wallet address is required")
	}

	if request.Network == "" {
		request.Network = "solana-devnet"
	}

	return s.repository.Create(
		ctx,
		userID,
		request.Name,
		request.Address,
		request.Network,
	)
}

func (s *walletService) Get(
	ctx context.Context,
	id string,
	userID string,
) (*models.Wallet, error) {

	return s.repository.FindByID(
		ctx,
		id,
		userID,
	)
}

func (s *walletService) List(
	ctx context.Context,
	userID string,
) ([]models.Wallet, error) {

	return s.repository.FindByUserID(
		ctx,
		userID,
	)
}
