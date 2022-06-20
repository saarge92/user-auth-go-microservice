package transformer

import (
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/pkg/protobuf/core"
)

func WalletsDtoToGrpc(walletsDto []dto.WalletCurrencyDto) *core.WalletsResponse {
	walletsResponse := make([]*core.Wallet, 0, len(walletsDto))
	for _, walletDto := range walletsDto {
		walletInstance := walletDto.Wallet
		walletResponseElement := &core.Wallet{
			ExternalId: walletInstance.ExternalID,
			Currency:   walletDto.Currency.Code,
			Balance:    walletInstance.Balance.String(),
			IsDefault:  walletInstance.IsDefault,
		}
		walletsResponse = append(walletsResponse, walletResponseElement)
	}
	return &core.WalletsResponse{
		Wallets: walletsResponse,
	}
}
