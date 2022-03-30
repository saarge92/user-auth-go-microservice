package db

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/errors"
)

type TransactionHandlerDB struct {
	dbConn *sqlx.DB
}

func NewTransactionHandler(dbConn *sqlx.DB) *TransactionHandlerDB {
	return &TransactionHandlerDB{
		dbConn: dbConn,
	}
}

func (t *TransactionHandlerDB) Create() *sqlx.Tx {
	return t.dbConn.MustBegin()
}

func (t *TransactionHandlerDB) Commit(tx *sqlx.Tx) error {
	if e := tx.Commit(); e != nil {
		return errors.DatabaseError(e)
	}
	return nil
}

func (t *TransactionHandlerDB) Rollback(tx *sqlx.Tx) error {
	if e := tx.Rollback(); e != nil {
		return errors.DatabaseError(e)
	}
	return nil
}
