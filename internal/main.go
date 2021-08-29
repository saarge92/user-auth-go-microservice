package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	specifyLogger()
}

func specifyLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.StandardLogger()
	log.SetOutput(os.Stdout)
}
