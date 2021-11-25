package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/wallet/entities"
	customErrors "go-user-microservice/internal/pkg/errors"
	sharedRepositories "go-user-microservice/internal/pkg/repositories"
	"time"
)

type WalletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
	wallet.UpdatedAt = time.Now()
	wallet.CreatedAt = time.Now()
	wallet.ExternalID = uuid.New().String()
	query := `INSERT INTO wallets(
                    currency_id, user_id, balance, 
                    external_id, is_default, created_at, updated_at)
				VALUES (:currency_id, :user_id, :balance,
				       :external_id, :is_default, :created_at, :updated_at)`
	var result sql.Result
	var dbError error
	tx := sharedRepositories.GetDBTransaction(ctx)
	if tx != nil {
		result, dbError = tx.NamedExec(query, wallet)
	} else {
		result, dbError = r.db.NamedExec(query, wallet)
	}
	if dbError != nil {
		return customErrors.DatabaseError(dbError)
	}
	wallet.ID = uint64(sharedRepositories.LastInsertID(result))
	return nil
}

func (r *WalletRepository) Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error) {
	query := `SELECT * FROM wallets WHERE user_id = ? AND currency_id = ?`
	wallet := &entities.Wallet{}

	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Get(wallet, query, userID, currencyID)
	} else {
		dbError = r.db.Get(wallet, query, userID, currencyID)
	}
	if dbError != nil {
		if dbError == sql.ErrNoRows {
			return false, nil
		}
		return false, dbError
	}
	return true, nil
}

func (r *WalletRepository) ByUserAndDefault(
	ctx context.Context,
	userID uint64,
	isDefault bool,
) (*entities.Wallet, error) {
	query := `SELECT * FROM wallets WHERE user_id = ? AND is_default = ?`
	wallet := &entities.Wallet{}
	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		dbError = tx.Get(wallet, query, userID, isDefault)
	} else {
		dbError = r.db.Get(wallet, query, userID, isDefault)
	}
	if dbError != nil {
		if dbError == sql.ErrNoRows {
			return nil, nil
		}
		return nil, customErrors.DatabaseError(dbError)
	}
	return wallet, nil
}

func (r *WalletRepository) UpdateStatusByUserID(ctx context.Context, userID uint64, isDefault bool) error {
	query := `UPDATE wallets SET is_default = ? WHERE user_id = ?`
	tx := sharedRepositories.GetDBTransaction(ctx)
	var dbError error
	if tx != nil {
		_, dbError = tx.Exec(query, isDefault, userID)
	} else {
		_, dbError = r.db.Exec(query, isDefault, userID)
	}
	if dbError != nil {
		return dbError
	}
	return nil
}
