package providers

import (
	"github.com/jmoiron/sqlx"
	"go-user-microservice/internal/config"
	"log"
)

type ConnectionProvider struct {
	coreConnection *sqlx.DB
}

func NewConnectionProvider(config *config.Config) *ConnectionProvider {
	coreConn, e := sqlx.Connect("mysql", config.CoreDatabaseURL)
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

func (c *ConnectionProvider) SetCoreConnection(conn *sqlx.DB) {
	c.coreConnection = conn
}
