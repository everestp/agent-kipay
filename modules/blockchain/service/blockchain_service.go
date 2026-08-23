package service

import (
	"context"
	"errors"
	"strings"

	"github.com/everest/bheri/modules/blockchain/client"
	"github.com/everest/bheri/modules/blockchain/dto"
)

type BlockchainService interface {
	VerifyTransaction(
		ctx context.Context,
		req dto.VerifyTransactionRequest,
	) (*dto.VerifyTransactionResponse, error)
}

type blockchainService struct {
	solanaClient *client.SolanaClient
}

func NewBlockchainService(
	solanaClient *client.SolanaClient,
) BlockchainService {
	return &blockchainService{
		solanaClient: solanaClient,
	}
}

func (s *blockchainService) VerifyTransaction(
	ctx context.Context,
	req dto.VerifyTransactionRequest,
) (*dto.VerifyTransactionResponse, error) {

	if req.TxHash == "" {
		return nil, errors.New("transaction hash is required")
	}

	if req.Network == "" {
		return nil, errors.New("network is required")
	}

	if req.Network != "solana" &&
		req.Network != "solana-devnet" {
		return nil, errors.New("unsupported network")
	}

	tx, err := s.solanaClient.GetTransaction(
		ctx,
		req.TxHash,
	)

	if err != nil {
		return nil, err
	}

	if tx.Meta == nil {
		return nil, errors.New(
			"transaction metadata unavailable",
		)
	}

	if tx.Meta.Err != nil {
		return &dto.VerifyTransactionResponse{
			Valid:       false,
			Confirmed:   false,
			TxHash:      req.TxHash,
			Network:     req.Network,
			BlockNumber: tx.Slot,
			Message:     "transaction failed on chain",
		}, nil
	}

	/*
		At this layer we know that:

		1. transaction exists
		2. transaction has confirmed metadata
		3. transaction did not fail

		Token amount / receiver verification should be
		added by parsing the transaction instructions
		or token balance changes.
	*/

	receiver := ""

	if len(tx.Transaction.Message.AccountKeys) > 0 {
		receiver =
			tx.Transaction.Message.AccountKeys[0].Pubkey
	}

	valid := true

	if req.ExpectedReceiver != "" &&
		receiver != req.ExpectedReceiver {
		valid = false
	}

	message := "transaction verified"

	if !valid {
		message = "receiver mismatch"
	}

	return &dto.VerifyTransactionResponse{
		Valid:       valid,
		Confirmed:   true,
		TxHash:      req.TxHash,
		Network:     req.Network,
		Receiver:    receiver,
		BlockNumber: tx.Slot,
		Message:     message,
	}, nil
}

func normalizeNetwork(network string) string {
	return strings.ToLower(
		strings.TrimSpace(network),
	)
}
