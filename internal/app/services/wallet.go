package services

import (
	"github.com/shopspring/decimal"
	repositories2 "go-user-microservice/internal/app/domain/repositories"
	entites2 "go-user-microservice/internal/app/entites"
	errorlists2 "go-user-microservice/internal/app/errorlists"
	forms2 "go-user-microservice/internal/app/forms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   repositories2.WalletRepositoryInterface
	userRepository     repositories2.UserRepositoryInterface
	currencyRepository repositories2.CurrencyRepositoryInterface
}

func NewWalletService(
	walletRepository repositories2.WalletRepositoryInterface,
	userRepository repositories2.UserRepositoryInterface,
	currencyRepository repositories2.CurrencyRepositoryInterface,
) *WalletService {
	return &WalletService{
		walletRepository:   walletRepository,
		userRepository:     userRepository,
		currencyRepository: currencyRepository,
	}
}

func (s *WalletService) Create(form *forms2.WalletCreateForm) (*entites2.Wallet, error) {
	user, e := s.userRepository.UserByID(form.UserID)
	if e != nil {
		return nil, e
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, errorlists2.UserNotFound)
	}
	currency, e := s.currencyRepository.GetByCode(form.Code)
	if e != nil {
		return nil, e
	}
	if currency == nil {
		return nil, status.Error(codes.NotFound, errorlists2.CurrencyNotFound)
	}
	balance := decimal.NewFromInt(0)
	wallet := &entites2.Wallet{
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
