package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/providers"
	"os"
	"path"
	"runtime"
)

type Server struct {
	userGrpcServer   *user.GrpcUserServer
	walletGrpcServer *wallet.GrpcWalletServer
	cardGrpcServer   *card.GrpcServerCard
}

func NewServer() *Server {
	mainServer := &Server{}
	return mainServer
}
func (s *Server) InitConfig() error {
	var _, filename, _, _ = runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	e := os.Chdir(dir)
	if e != nil {
		panic(e)
	}
	if e := godotenv.Load(".env"); e != nil {
		return e
	}
	return nil
}

func (s *Server) InitApp() error {
	appConfig := config.NewConfig()
	dbConnectionProvider := providers.NewConnectionProvider(appConfig)
	repositoryProvider := providers.NewRepositoryProvider(dbConnectionProvider)
	stripeClientProvider := providers.NewClientStripeProvider(appConfig)
	serviceProvider := providers.NewServiceProvider(
		appConfig,
		repositoryProvider,
		dbConnectionProvider,
		stripeClientProvider,
	)
	fmt.Println(serviceProvider)
	return nil
}
