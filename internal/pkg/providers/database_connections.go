package providers

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/pkg/config"
	"log"
)

type DatabaseConnectionProvider struct {
	coreConnection *sqlx.DB
}

func NewConnectionProvider(
	config *config.Config,
) *DatabaseConnectionProvider {
	coreConn, e := sqlx.Open(config.DatabaseDriver, config.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	return &DatabaseConnectionProvider{
		coreConnection: coreConn,
	}
}

func (c *DatabaseConnectionProvider) GetCoreConnection() *sqlx.DB {
	return c.coreConnection
}
