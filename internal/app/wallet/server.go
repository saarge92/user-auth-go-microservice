package wallet

import (
	"context"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/domain/services"
	walletGrpcServer "go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcWalletServer struct {
	walletService services.WalletServiceInterface
}

func (s *GrpcWalletServer) Wallets(
	ctx context.Context,
	empty *emptypb.Empty,
) (*walletGrpcServer.WalletsResponse, error) {
	return nil, nil
}

func NewWalletGrpcServer(
	walletService services.WalletServiceInterface,
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
