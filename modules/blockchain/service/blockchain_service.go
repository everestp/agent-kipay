package service

import (
	"context"

	"github.com/everest/bheri/modules/blockchain/dto"
)

type BlockchainService interface {

	Transfer(
		ctx context.Context,
		request dto.TransferRequest,
	) (*dto.TransferResponse, error)

	VerifyTransaction(
		ctx context.Context,
		request dto.VerifyTransactionRequest,
	) (*dto.VerifyTransactionResponse, error)
}
