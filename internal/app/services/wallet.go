package services

import (
	"context"
	"fmt"
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
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return s.initWallet(ctx, form)
	}
}

func (s *WalletService) initWallet(
	ctx context.Context,
	form *forms.WalletCreateForm) (*entites.Wallet, error) {
	user, currency, e := s.checkCreateWalletData(ctx, form)
	if e != nil {
		return nil, e
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

func (s *WalletService) checkCreateWalletData(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (*entites.User, *entites.Currency, error) {
	var userID uint64
	var ok bool
	userIDData := ctx.Value(dictionary.UserID)
	if userIDData == nil {
		return nil, nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	if userID, ok = ctx.Value(dictionary.UserID).(uint64); !ok {
		return nil, nil, status.Error(codes.Internal, fmt.Sprintf(errorlists.ConvertError, "user_id"))
	}
	user, e := s.userRepository.UserByID(userID)
	if e != nil {
		return nil, nil, e
	}
	if user == nil {
		return nil, nil, status.Error(codes.NotFound, errorlists.UserNotFound)
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
