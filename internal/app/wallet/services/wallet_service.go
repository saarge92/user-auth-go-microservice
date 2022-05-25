package services

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/entities"
	walletDomain "go-user-microservice/internal/app/wallet/domain"
	"go-user-microservice/internal/app/wallet/dto"
	walletEntities "go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/dictionary"
	repositoryInterface "go-user-microservice/internal/pkg/domain/repositories"
	sharedEntities "go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   walletDomain.WalletRepository
	userRepository     domain.UserRepository
	currencyRepository repositoryInterface.CurrencyRepository
}

func NewWalletService(
	walletRepository walletDomain.WalletRepository,
	userRepository domain.UserRepository,
	currencyRepository repositoryInterface.CurrencyRepository,
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
) (wallet *walletEntities.Wallet, e error) {
	user, currency, e := s.checkCreateWalletData(ctx, form)
	if e != nil {
		return nil, e
	}
	balance := decimal.NewFromInt(0)
	wallet = &walletEntities.Wallet{
		UserID:     user.ID,
		Balance:    balance,
		CurrencyID: currency.ID,
		IsDefault:  form.IsDefault,
	}
	if form.IsDefault {
		defaultWallet, e := s.walletRepository.ByUserAndDefault(ctx, user.ID, true)
		if e != nil {
			return nil, e
		}
		if defaultWallet != nil {
			e := s.walletRepository.UpdateStatusByUserID(ctx, user.ID, false)
			if e != nil {
				return nil, e
			}
		}
	}

	if e = s.walletRepository.Create(ctx, wallet); e != nil {
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
	currency, e := s.currencyRepository.GetByCode(ctx, form.Code)
	if e != nil {
		return nil, nil, e
	}
	if currency == nil {
		return nil, nil, status.Error(codes.NotFound, errorlists.CurrencyNotFound)
	}
	walletUserExist, e := s.walletRepository.Exist(ctx, user.ID, currency.ID)
	if e != nil {
		return nil, nil, e
	}
	if walletUserExist {
		return nil, nil, status.Error(codes.AlreadyExists, errorlists.UserWalletAlreadyExist)
	}
	return user, currency, nil
}

func (s *WalletService) Wallets(ctx context.Context) ([]dto.WalletCurrencyDto, error) {
	var user *entities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*entities.User); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}
	wallets, e := s.walletRepository.ListByUserID(ctx, user.ID)
	if e != nil {
		return nil, e
	}
	return wallets, nil
}
