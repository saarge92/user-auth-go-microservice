package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/dictionary"
	repositoryInterface "go-user-microservice/internal/pkg/domain/repositories"
	sharedEntities "go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   repositoryInterface.WalletRepositoryInterface
	userRepository     repositoryInterface.UserRepositoryInterface
	currencyRepository repositoryInterface.CurrencyRepositoryInterface
}

func NewWalletService(
	walletRepository repositoryInterface.WalletRepositoryInterface,
	userRepository repositoryInterface.UserRepositoryInterface,
	currencyRepository repositoryInterface.CurrencyRepositoryInterface,
) *WalletService {
	return &WalletService{
		walletRepository:   walletRepository,
		userRepository:     userRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *WalletService) Create(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (*sharedEntities.Wallet, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return s.initWallet(ctx, form)
	}
}

func (s *WalletService) initWallet(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (*sharedEntities.Wallet, error) {
	user, currency, e := s.checkCreateWalletData(ctx, form)
	if e != nil {
		return nil, e
	}
	balance := decimal.NewFromInt(0)
	wallet := &sharedEntities.Wallet{
		UserID:     user.ID,
		Balance:    balance,
		CurrencyID: currency.ID,
	}
	e = s.walletRepository.Create(wallet)
	if e != nil {
		return nil, e
	}
	return wallet, nil
}

func (s *WalletService) checkCreateWalletData(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (*entities.User, *sharedEntities.Currency, error) {
	var user *entities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*entities.User); !ok {
		return nil, nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	currency, e := s.currencyRepository.GetByCode(form.Code)
	if e != nil {
		return nil, nil, e
	}
	if currency == nil {
		return nil, nil, status.Error(codes.NotFound, errorlists.CurrencyNotFound)
	}
	walletUserExist, e := s.walletRepository.Exist(user.ID, currency.ID)
	if e != nil {
		return nil, nil, e
	}
	if walletUserExist {
		return nil, nil, status.Error(codes.AlreadyExists, errorlists.UserWalletAlreadyExist)
	}
	return user, currency, nil
}
