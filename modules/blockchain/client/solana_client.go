package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

)

type SolanaClient struct {
	rpcURL string
	client *http.Client
}

func NewSolanaClient(rpcURL string) *SolanaClient {
	return &SolanaClient{
		rpcURL: rpcURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type transactionResult struct {
	BlockTime   *int64 `json:"blockTime"`
	Slot        uint64 `json:"slot"`
	Meta        *struct {
		Err interface{} `json:"err"`
	} `json:"meta"`
	Transaction struct {
		Message struct {
			AccountKeys []struct {
				Pubkey string `json:"pubkey"`
			} `json:"accountKeys"`
		} `json:"message"`
	} `json:"transaction"`
}

func (c *SolanaClient) GetTransaction(
	ctx context.Context,
	txHash string,
) (*transactionResult, error) {

	reqBody := rpcRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "getTransaction",
		Params: []interface{}{
			txHash,
			map[string]interface{}{
				"encoding":                       "json",
				"commitment":                     "confirmed",
				"maxSupportedTransactionVersion": 0,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.rpcURL,
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"solana rpc returned status %d",
			resp.StatusCode,
		)
	}

	var result rpcResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, errors.New(result.Error.Message)
	}

	if string(result.Result) == "null" {
		return nil, errors.New("transaction not found")
	}

	var tx transactionResult

	if err := json.Unmarshal(
		result.Result,
		&tx,
	); err != nil {
		return nil, err
	}

	return &tx, nil
}
