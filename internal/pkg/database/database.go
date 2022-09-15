package database

import (
	"context"
	"database/sql"
)

type database struct {
	databaseConnection *sql.DB
}

func NewDatabase(databaseConnection *sql.DB) Database {
	return &database{
		databaseConnection: databaseConnection,
	}
}

func (d *database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return getDatabaseConnection(ctx, d.databaseConnection).QueryRowContext(ctx, query, args...)
}

func (d *database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return getDatabaseConnection(ctx, d.databaseConnection).ExecContext(ctx, query, args...)
}

func (d *database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return getDatabaseConnection(ctx, d.databaseConnection).QueryContext(ctx, query, args...)
}

func (d *database) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return getDatabaseConnection(ctx, d.databaseConnection).PrepareContext(ctx, query)
}

func getDatabaseConnection(ctx context.Context, dbConnection *sql.DB) Database {
	var returnDatabase Database
	currentTransaction, ok := ctx.Value(CurrentTransaction).(*sql.Tx)
	if ok {
		returnDatabase = currentTransaction
	} else {
		returnDatabase = dbConnection
	}
	return returnDatabase
}
