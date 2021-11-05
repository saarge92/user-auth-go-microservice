package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
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
