package wallet

import (
	"context"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/app/wallet/transformer"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcWalletServer struct {
	walletService      services.WalletService
	transactionHandler *db.TransactionHandlerDB
}

func NewWalletGrpcServer(
	walletService services.WalletService,
	transactionHandler *db.TransactionHandlerDB,
) *GrpcWalletServer {
	return &GrpcWalletServer{
		walletService:      walletService,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcWalletServer) CreateWallet(
	ctx context.Context,
	message *core.CreateWalletRequest,
) (resp *core.CreateWalletResponse, e error) {
	walletCreateForm := forms.NewWalletCreateForm(message)
	if e = walletCreateForm.Validate(); e != nil {
		return nil, e
	}

	ctx, handleFunc := db.MakeConnectionContext(ctx, s.transactionHandler)
	defer func() {
		e = handleFunc(e)
	}()

	walletEntity, e := s.walletService.Create(ctx, walletCreateForm)
	if e != nil {
		return nil, e
	}
	return &core.CreateWalletResponse{
		ExternalId: walletEntity.ExternalID,
	}, nil
}

func (s *GrpcWalletServer) MyWallets(
	ctx context.Context,
	_ *core.MyWalletsRequest,
) (*core.WalletsResponse, error) {
	wallets, e := s.walletService.Wallets(ctx)
	if e != nil {
		return nil, e
	}
	return transformer.WalletsDtoToGrpc(wallets), nil
}
