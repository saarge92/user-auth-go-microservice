package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/domain"
	"go-user-microservice/internal/pkg/errors"
)

type TransactionKey int

const (
	CurrentTransaction TransactionKey = iota
)

func GetDBConnection(ctx context.Context, dbConn *sqlx.DB) domain.SQLDb {
	var ret domain.SQLDb
	tx, ok := ctx.Value(CurrentTransaction).(*sqlx.Tx)
	if ok {
		ret = tx
	} else {
		ret = dbConn
	}
	return ret
}

func MakeConnectionContext(
	ctx context.Context,
	transactionHandler *TransactionHandlerDB,
) (context.Context, func(e error) error) {
	tx := transactionHandler.Create()
	newCtx := context.WithValue(ctx, CurrentTransaction, tx)

	return newCtx, func(e error) error {
		return HandleTransaction(tx, e)
	}
}

func HandleTransaction(tx *sqlx.Tx, functionError error) error {
	if functionError != nil {
		if rollErr := tx.Rollback(); rollErr != nil {
			return errors.DatabaseError(rollErr)
		}

		return errors.DatabaseError(functionError)
	}

	if e := tx.Commit(); e != nil {
		return errors.DatabaseError(e)
	}

	return nil
}

func LastInsertID(result sql.Result) int32 {
	id, _ := result.LastInsertId()
	return int32(id)
}
