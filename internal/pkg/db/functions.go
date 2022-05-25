package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	sharedErrors "go-user-microservice/internal/pkg/errors"
	"go-user-microservice/internal/pkg/filter"
)

type TransactionKey int

const (
	CurrentTransaction TransactionKey = iota
)

func LastInsertID(result sql.Result) int32 {
	id, _ := result.LastInsertId()
	return int32(id)
}

func AddPagination(query string, pagination filter.Pagination) string {
	if pagination.Page == 0 {
		return query + fmt.Sprintf(" LIMIT %d", pagination.PerPage)
	}
	return query + fmt.Sprintf(" LIMIT %d, %d", (pagination.Page-1)*pagination.PerPage, pagination.PerPage)
}

func HandleTransaction(tx driver.Tx, functionError error) error {
	if functionError != nil {
		if rollErr := tx.Rollback(); rollErr != nil {
			return sharedErrors.DatabaseError(rollErr)
		}

		return sharedErrors.DatabaseError(functionError)
	}

	if e := tx.Commit(); e != nil {
		return sharedErrors.DatabaseError(e)
	}

	return nil
}
