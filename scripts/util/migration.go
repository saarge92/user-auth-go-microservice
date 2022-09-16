package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/scripts/migrations"
)

func main() {
	envFileName := flag.String("env-file", ".env", "environment file name")
	flag.Parse()
	if err := godotenv.Load(*envFileName); err != nil {
		panic(err)
	}
	appConfig := config.NewConfig()
	fmt.Println(appConfig.CoreDatabaseURL)
	migrations.Migrate(appConfig)
}
