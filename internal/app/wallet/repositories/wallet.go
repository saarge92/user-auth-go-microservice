package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Knetic/go-namedParameterQuery"
	"github.com/blockloop/scan"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/internal/app/wallet/entities"
	walletErrors "go-user-microservice/internal/app/wallet/errors"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/errorlists"
	customErrors "go-user-microservice/internal/pkg/errors"
	"google.golang.org/grpc/codes"
	"time"
)

type WalletRepository struct {
	databaseInstance database.Database
}

func NewWalletRepository(databaseInstance database.Database) *WalletRepository {
	return &WalletRepository{
		databaseInstance: databaseInstance,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
	wallet.UpdatedAt = time.Now()
	wallet.CreatedAt = time.Now()
	wallet.ExternalID = uuid.New().String()
	query := `INSERT INTO wallets(
                    currency_id, user_id, balance, 
                    external_id, is_default, created_at, updated_at)
				VALUES (:currencyId, :userId, :balance,
				       :externalId, :isDefault, :createdAt, :updatedAt)`

	queryNamed := namedParameterQuery.NewNamedParameterQuery(query)
	insertParams := map[string]interface{}{
		"currencyId": wallet.CurrencyID,
		"userId":     wallet.UserID,
		"balance":    wallet.Balance,
		"externalId": wallet.ExternalID,
		"isDefault":  wallet.IsDefault,
		"createdAt":  wallet.CreatedAt,
		"updatedAt":  wallet.UpdatedAt,
	}
	queryNamed.SetValuesFromMap(insertParams)
	result, dbError := r.databaseInstance.ExecContext(ctx, queryNamed.GetParsedQuery(), queryNamed.GetParsedParameters()...)

	if dbError != nil {
		return dbError
	}

	wallet.ID = uint64(db.LastInsertID(result))
	return nil
}

func (r *WalletRepository) Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error) {
	query := `SELECT COUNT(*) > 0 FROM wallets WHERE user_id = ? AND currency_id = ?`
	var exist bool

	if e := r.databaseInstance.QueryRowContext(ctx, query, userID, currencyID).Scan(&exist); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return false, nil
		}
		return false, customErrors.DatabaseError(e)
	}

	return exist, nil
}

func (r *WalletRepository) DefaultWalletByUser(
	ctx context.Context,
	userID uint64,
) (*entities.Wallet, error) {
	query := `SELECT * FROM wallets WHERE user_id = ? AND is_default = ?`
	wallet := &entities.Wallet{}

	walletRows, walletError := r.databaseInstance.QueryContext(ctx, query, userID, true)
	if walletError != nil {
		return nil, customErrors.DatabaseError(walletError)
	}

	if e := scan.Row(wallet, walletRows); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, walletErrors.WalletNotFoundErr
		}
		return nil, customErrors.DatabaseError(e)
	}

	return wallet, nil
}

func (r *WalletRepository) SetAsDefaultForUserWallet(ctx context.Context, userID uint64, isDefault bool) error {
	query := `UPDATE wallets SET is_default = ? WHERE user_id = ?`

	result, dbError := r.databaseInstance.ExecContext(ctx, query, isDefault, userID)
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	if countAffected, _ := result.LastInsertId(); countAffected == 0 {
		return walletErrors.WalletNotFoundErr
	}

	return nil
}

func (r *WalletRepository) ListByUserID(ctx context.Context, userID uint64) ([]dto.WalletCurrencyDto, error) {
	query := `SELECT 
				wallets.id "wallet.id",
       			wallets.external_id "wallet.external_id",
       			wallets.balance "wallet.balance",
       			wallets.is_default "wallet.balance",
       			currencies.code "currency.code"
				FROM wallets 
				INNER JOIN currencies ON wallets.currency_id = currencies.id
				WHERE wallets.user_id = ?
			`
	var walletsCurrencies []dto.WalletCurrencyDto

	walletRows, walletError := r.databaseInstance.QueryContext(ctx, query, userID)
	if walletError != nil {
		return nil, customErrors.DatabaseError(walletError)
	}

	if e := scan.Rows(&walletsCurrencies, walletRows); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, customErrors.DatabaseError(e)
	}

	return walletsCurrencies, nil
}

func (r *WalletRepository) OneByExternalIDAndUserID(
	ctx context.Context,
	externalID string,
	userID uint64,
) (*dto.WalletCurrencyDto, error) {
	query := `SELECT 
				wallets.id "wallet.id",
       			wallets.external_id "wallet.external_id",
       			wallets.balance "wallet.balance",
       			wallets.is_default "wallet.balance",
       			currencies.code "currency.code"
				FROM wallets 
				INNER JOIN currencies ON wallets.currency_id = currencies.id
				WHERE wallets.user_id = ? AND wallets.external_id = ?
			`
	walletWithCurrencyDto := new(dto.WalletCurrencyDto)

	walletRow, walletError := r.databaseInstance.QueryContext(ctx, query, userID, externalID)
	if walletError != nil {
		return nil, customErrors.DatabaseError(walletError)
	}

	if e := scan.Row(walletWithCurrencyDto, walletRow); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return nil, walletErrors.WalletNotFoundErr
		}
		return nil, customErrors.DatabaseError(e)
	}

	return walletWithCurrencyDto, nil
}

func (r *WalletRepository) IncreaseBalanceByID(ctx context.Context, id uint64, amount decimal.Decimal) error {
	updatedAt := time.Now()
	query := `UPDATE wallets SET balance = balance + ?,
				updated_at = ? WHERE id = ?`

	var dbError error
	var result sql.Result
	result, dbError = r.databaseInstance.ExecContext(ctx, query, amount, updatedAt, id)

	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return customErrors.CustomDatabaseError(codes.NotFound, errorlists.WalletNotFound)
	}

	return nil
}
