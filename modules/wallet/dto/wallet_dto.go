// modules/wallet/dto/wallet_dto.go

package dto

type CreateWalletRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Network string `json:"network"`
}
