package dto

type VerifyTransactionRequest struct {
	TxHash          string
	Network         string
	ExpectedAmount  float64
	ExpectedAsset   string
	ExpectedReceiver string
}

type VerifyTransactionResponse struct {
	Valid         bool
	Confirmed     bool
	TxHash        string
	Network       string
	Amount        float64
	Asset         string
	Sender        string
	Receiver      string
	BlockNumber   uint64
	Message       string
}
