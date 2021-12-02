package providers

import "github.com/jmoiron/sqlx"

type DatabaseConnectionProvider interface {
	GetCoreConnection() *sqlx.DB
}
