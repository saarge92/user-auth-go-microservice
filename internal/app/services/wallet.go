package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/forms"
	"go-user-microservice/internal/app/middlewares/dictionary"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   repositories.WalletRepositoryInterface
	userRepository     repositories.UserRepositoryInterface
	currencyRepository repositories.CurrencyRepositoryInterface
}

func NewWalletService(
	walletRepository repositories.WalletRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
	currencyRepository repositories.CurrencyRepositoryInterface,
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
) (*entites.Wallet, error) {
	var userID uint64
	var ok bool
	if userID, ok = ctx.Value(dictionary.UserID).(uint64); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	user, e := s.userRepository.UserByID(userID)
	if e != nil {
		return nil, e
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, errorlists.UserNotFound)
	}
	currency, e := s.currencyRepository.GetByCode(form.Code)
	if e != nil {
		return nil, e
	}
	if currency == nil {
		return nil, status.Error(codes.NotFound, errorlists.CurrencyNotFound)
	}
	walletUserExist, e := s.walletRepository.Exist(user.ID, currency.ID)
	if e != nil {
		return nil, e
	}
	if walletUserExist {
		return nil, status.Error(codes.AlreadyExists, errorlists.UserWalletAlreadyExist)
	}
	balance := decimal.NewFromInt(0)
	wallet := &entites.Wallet{
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
