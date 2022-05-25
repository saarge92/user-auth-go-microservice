package db

import (
	"context"
	"database/sql"
)

type TransactionHandlerDB struct {
	dbConn *sql.DB
}

func NewTransactionHandler(dbConn *sql.DB) *TransactionHandlerDB {
	return &TransactionHandlerDB{
		dbConn: dbConn,
	}
}

func (t *TransactionHandlerDB) Create(ctx context.Context, options *sql.TxOptions) (context.Context, *sql.Tx, error) {
	tx, e := t.dbConn.BeginTx(ctx, options)
	if e != nil {
		return ctx, nil, e
	}
	return context.WithValue(ctx, CurrentTransaction, tx), tx, nil
}

func (t *TransactionHandlerDB) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (t *TransactionHandlerDB) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
