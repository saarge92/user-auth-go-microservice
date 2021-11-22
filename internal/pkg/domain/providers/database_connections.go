package providers

import "github.com/jmoiron/sqlx"

type DatabaseConnectionProviderInterface interface {
	GetCoreConnection() *sqlx.DB
}
