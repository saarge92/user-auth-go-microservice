package providers

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/config"
	"log"
)

type ConnectionProvider struct {
	coreConnection *sqlx.DB
}

func NewConnectionProvider(
	config *config.Config,
	driverDB string,
) *ConnectionProvider {
	coreConn, e := sqlx.Connect(driverDB, config.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	return &ConnectionProvider{
		coreConnection: coreConn,
	}
}

func NewConnectionProviderForConnection(connection *sqlx.DB) *ConnectionProvider {
	return &ConnectionProvider{coreConnection: connection}
}

func (c *ConnectionProvider) GetCoreConnection() *sqlx.DB {
	return c.coreConnection
}
