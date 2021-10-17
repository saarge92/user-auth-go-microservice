package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/app/entites"
	customErrors "go-user-microservice/internal/app/errors"
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

func (r *WalletRepository) Create(wallet *entites.Wallet) error {
	wallet.UpdatedAt = time.Now()
	wallet.CreatedAt = time.Now()
	query := `INSERT INTO wallets(currency_id, user_id, balance, created_at, updated_at)
				VALUES (:currency_id, :user_id, :balance, :created_at, :updated_at)`
	result, e := r.db.NamedExec(query, wallet)
	if e != nil {
		return customErrors.DatabaseError(e)
	}
	wallet.ID = uint64(lastInsertID(result))
	return nil
}

func (r *WalletRepository) Exist(userID uint64, currencyID uint32) (bool, error) {
	query := `SELECT * from wallets WHERE user_id = ? AND currency_id = ?`
	wallet := &entites.Wallet{}
	e := r.db.Get(wallet, query, userID, currencyID)
	if e != nil {
		if e == sql.ErrNoRows {
			return false, nil
		}
		return false, e
	}
	return true, nil
}
