package migrations

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/pkg/config"
	"os"
)

func Migrate(configSetting *config.Config) {
	if configSetting == nil {
		configSetting = config.NewConfig()
	}
	db, e := sql.Open("mysql", configSetting.CoreDatabaseURL)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	driver, e := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})
	if e != nil {
		log.Fatal(e)
	}
	migrationDir := getScriptDirectory()
	migration, e := migrate.NewWithDatabaseInstance(
		"file://"+migrationDir,
		"mysql",
		driver,
	)
	if e != nil {
		log.Fatal(e)
	}

	if e := migration.Up(); e != nil {
		log.Println(e)
	}
	migration.Close()
}

func getScriptDirectory() string {
	d, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	return d + "/scripts/migrations/sql/"
}
