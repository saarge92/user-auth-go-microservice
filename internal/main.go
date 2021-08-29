package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	specifyLogger()
	if err := godotenv.Load(".env"); err != nil {
		log.Debug(".env file not found, use system environment")
	}
	server := NewServer()
	if e := server.Start(); e != nil {
		log.Fatal(e)
	}
	fmt.Println("Microservice has been launched")
}

func specifyLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.StandardLogger()
	log.SetOutput(os.Stdout)
}
