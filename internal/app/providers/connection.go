package providers

import (
	"github.com/jmoiron/sqlx"
	config2 "go-user-microservice/internal/app/config"
	"log"
)

type ConnectionProvider struct {
	coreConnection *sqlx.DB
}

func NewConnectionProvider(
	config *config2.Config,
) *ConnectionProvider {
	coreConn, e := sqlx.Open(config.DatabaseDriver, config.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	return &ConnectionProvider{
		coreConnection: coreConn,
	}
}

func (c *ConnectionProvider) GetCoreConnection() *sqlx.DB {
	return c.coreConnection
}
