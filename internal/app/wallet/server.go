package wallet

import (
	"context"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/domain/services"
	walletGrpcServer "go-user-microservice/pkg/protobuf/wallet"
)

type GrpcWalletServer struct {
	walletService services.WalletService
}

func NewWalletGrpcServer(
	walletService services.WalletService,
) *GrpcWalletServer {
	return &GrpcWalletServer{
		walletService: walletService,
	}
}

func (s *GrpcWalletServer) CreateWallet(
	ctx context.Context,
	message *walletGrpcServer.CreateWalletRequest,
) (*walletGrpcServer.CreateWalletResponse, error) {
	walletCreateForm := forms.NewWalletCreateForm(message)
	if e := walletCreateForm.Validate(); e != nil {
		return nil, e
	}
	walletEntity, e := s.walletService.Create(ctx, walletCreateForm)
	if e != nil {
		return nil, e
	}
	return &walletGrpcServer.CreateWalletResponse{
		ExternalId: walletEntity.ExternalID,
	}, nil
}

func (s *GrpcWalletServer) MyWallets(
	ctx context.Context,
	_ *walletGrpcServer.MyWalletsRequest,
) (*walletGrpcServer.WalletsResponse, error) {
	wallets, e := s.walletService.Wallets(ctx)
	if e != nil {
		return nil, e
	}
	return WalletsDtoToGrpc(wallets), nil
}
