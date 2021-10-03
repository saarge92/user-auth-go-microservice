package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	main2 "go-user-microservice/internal/app"
	"os"
)

func main() {
	specifyLogger()
	if err := godotenv.Load(".env"); err != nil {
		log.Debug(".env file not found, use system environment")
	}
	server := main2.NewServer()
	if e := server.Start(); e != nil {
		log.Fatal(e)
	}
}

func specifyLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.StandardLogger()
	log.SetOutput(os.Stdout)
}
