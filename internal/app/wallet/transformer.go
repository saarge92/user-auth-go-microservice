package wallet

import (
	"go-user-microservice/internal/app/wallet/dto"
	walletGrpcServer "go-user-microservice/pkg/protobuf/wallet"
)

func WalletsDtoToGrpc(walletsDto []dto.WalletCurrencyDto) *walletGrpcServer.WalletsResponse {
	walletsResponse := make([]*walletGrpcServer.Wallet, 0, len(walletsDto))
	for _, walletDto := range walletsDto {
		walletResponseElement := &walletGrpcServer.Wallet{
			ExternalId: walletDto.ExternalID,
			Currency:   walletDto.Code,
			Balance:    walletDto.Balance.String(),
			IsDefault:  walletDto.IsDefault,
		}
		walletsResponse = append(walletsResponse, walletResponseElement)
	}
	return &walletGrpcServer.WalletsResponse{
		Wallets: walletsResponse,
	}
}
