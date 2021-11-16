package services

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/dictionary"
	repositoryInterface "go-user-microservice/internal/pkg/domain/repositories"
	sharedEntities "go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/internal/pkg/repositories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletService struct {
	walletRepository   repositoryInterface.WalletRepositoryInterface
	userRepository     repositoryInterface.UserRepositoryInterface
	currencyRepository repositoryInterface.CurrencyRepositoryInterface
	coreDB             *sqlx.DB
}

func NewWalletService(
	walletRepository repositoryInterface.WalletRepositoryInterface,
	userRepository repositoryInterface.UserRepositoryInterface,
	currencyRepository repositoryInterface.CurrencyRepositoryInterface,
	coreDB *sqlx.DB,
) *WalletService {
	return &WalletService{
		walletRepository:   walletRepository,
		userRepository:     userRepository,
		currencyRepository: currencyRepository,
		coreDB:             coreDB,
	}
}

func (s *WalletService) Create(
	ctx context.Context,
	form *forms.WalletCreateForm,
) (wallet *sharedEntities.Wallet, e error) {
	tx := s.coreDB.MustBegin()
	defer func() {
		e = repositories.HandleTransaction(tx, e)
	}()
	newCtx := context.WithValue(ctx, repositories.CurrentTransaction, tx)
	user, currency, e := s.checkCreateWalletData(newCtx, form)
	if e != nil {
		return nil, e
	}
	balance := decimal.NewFromInt(0)
	wallet = &sharedEntities.Wallet{
		UserID:     user.ID,
		Balance:    balance,
		CurrencyID: currency.ID,
		IsDefault:  form.IsDefault,
	}
	if form.IsDefault {
		defaultWallet, e := s.walletRepository.ByUserAndDefault(newCtx, user.ID, true)
		if e != nil {
			return nil, e
		}
		if defaultWallet != nil {
			e := s.walletRepository.UpdateStatusByUserID(newCtx, user.ID, false)
			if e != nil {
				return nil, e
			}
		}
	}

	if e = s.walletRepository.Create(newCtx, wallet); e != nil {
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
