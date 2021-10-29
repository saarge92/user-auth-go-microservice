package server

import (
	"context"
	"go-user-microservice/internal/app/domain/services"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/forms"
	"go-user-microservice/internal/app/middlewares/dictionary"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WalletGrpcServer struct {
	walletService services.WalletServiceInterface
}

func NewWalletGrpcServer(
	walletService services.WalletServiceInterface,
) *WalletGrpcServer {
	return &WalletGrpcServer{
		walletService: walletService,
	}
}

func (s *WalletGrpcServer) CreateWallet(
	ctx context.Context,
	message *wallet.CreateWalletMessage,
) (*emptypb.Empty, error) {
	var userID uint64
	var ok bool
	if userID, ok = ctx.Value(dictionary.UserID).(uint64); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	walletCreateForm := forms.NewWalletCreateForm(message, userID)
	walletCreateForm.UserID = userID
	if e := walletCreateForm.Validate(); e != nil {
		return nil, e
	}
	_, e := s.walletService.Create(walletCreateForm)
	if e != nil {
		return nil, e
	}
	return &emptypb.Empty{}, nil
}
