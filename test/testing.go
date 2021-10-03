package test

import (
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app"
)

const (
	Inn   uint64 = 7707083893
	Login string = "user@mail.ru"
)

func CreateTestServer() *app.Server {
	if e := godotenv.Load(".env.test"); e != nil {
		panic(e)
	}
	server := app.NewServer()
	return server
}
