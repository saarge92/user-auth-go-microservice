package repositories

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/errors"
)

func LastInsertID(result sql.Result) int64 {
	id, _ := result.LastInsertId()
	return id
}

func GetDBTransaction(ctx context.Context) *sqlx.Tx {
	var dbTransaction *sqlx.Tx
	tx, ok := ctx.Value(db.CurrentTransaction).(*sqlx.Tx)
	if ok {
		dbTransaction = tx
	}
	return dbTransaction
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
