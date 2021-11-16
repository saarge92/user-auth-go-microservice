package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/errors"
)

func LastInsertID(result sql.Result) int32 {
	id, _ := result.LastInsertId()
	fmt.Println(id)
	return int32(id)
}

func GetDBTransaction(ctx context.Context) *sqlx.Tx {
	var dbTransaction *sqlx.Tx
	tx, ok := ctx.Value(CurrentTransaction).(*sqlx.Tx)
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
