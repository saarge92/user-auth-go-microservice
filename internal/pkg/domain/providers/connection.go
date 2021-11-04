package providers

import "github.com/jmoiron/sqlx"

type ConnectionProviderInterface interface {
	GetCoreConnection() *sqlx.DB
}
