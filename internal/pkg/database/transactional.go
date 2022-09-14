package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	log "github.com/sirupsen/logrus"
)

var lastConnection = 0

func TestDBConnection(dbURL string) (*sql.DB, func()) {
	lastConnection++
	connName := fmt.Sprintf("txdb_%d", lastConnection)
	log.Infof("Connection number: %d", lastConnection)

	txdb.Register(connName, "mysql", dbURL)

	databaseInstance, err := sql.Open(connName, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	databaseInstance.SetMaxIdleConns(1)
	return databaseInstance, func() {
		if e := databaseInstance.Close(); e != nil {
			log.Error(e)
		}
		log.Infof("Connection closed: %s", connName)
	}
}
