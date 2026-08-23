package client

import (
	"context"

	"github.com/everest/bheri/modules/blockchain/dto"
)

type SolanaClient interface {

	TransferUSDC(
		ctx context.Context,
		request dto.TransferRequest,
	) (*dto.TransferResponse, error)

	GetTransaction(
		ctx context.Context,
		txHash string,
	) (*dto.VerifyTransactionResponse, error)
}
