package wallet

import (
	"context"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcWalletServer struct {
	walletService services.WalletServiceInterface
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
	message *wallet.CreateWalletMessage,
) (*emptypb.Empty, error) {
	walletCreateForm := forms.NewWalletCreateForm(message)
	if e := walletCreateForm.Validate(); e != nil {
		return nil, e
	}
	_, e := s.walletService.Create(ctx, walletCreateForm)
	if e != nil {
		return nil, e
	}
	return &emptypb.Empty{}, nil
}
