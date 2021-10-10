package main

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app"
	"os"
)

func main() {
	specifyLogger()
	server := app.NewServer()
	if e := server.InitConfig(); e != nil {
		log.Fatal(e)
	}
	e := server.InitContainer("mysql")
	if e != nil {
		log.Fatal(e)
	}
	e = server.Start()
	if e != nil {
		log.Fatal(e)
	}
}

func specifyLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.StandardLogger()
	log.SetOutput(os.Stdout)
}
