package transformer

import (
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/pkg/protobuf/core"
)

func WalletsDtoToGrpc(walletsDto []dto.WalletCurrencyDto) *core.WalletsResponse {
	walletsResponse := make([]*core.Wallet, 0, len(walletsDto))
	for _, walletDto := range walletsDto {
		walletResponseElement := &core.Wallet{
			ExternalId: walletDto.ExternalID,
			Currency:   walletDto.Code,
			Balance:    walletDto.Balance.String(),
			IsDefault:  walletDto.IsDefault,
		}
		walletsResponse = append(walletsResponse, walletResponseElement)
	}
	return &core.WalletsResponse{
		Wallets: walletsResponse,
	}
}
