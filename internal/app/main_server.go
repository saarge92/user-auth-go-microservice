package app

import (
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/config"
	domainProviders "go-user-microservice/internal/pkg/domain/providers"
	"go-user-microservice/internal/pkg/providers"
	cardServer "go-user-microservice/pkg/protobuf/card"
	"go-user-microservice/pkg/protobuf/user_server"
	walletServer "go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
	"runtime"
)

type Server struct {
	mainConfig             *config.Config
	userGrpcServer         *user.GrpcUserServer
	walletGrpcServer       *wallet.GrpcWalletServer
	cardGrpcServer         *card.GrpcServerCard
	serviceProvider        domainProviders.ServiceProviderInterface
	grpcMiddlewareProvider *providers.GrpcMiddlewareProvider
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

func (s *Server) initApp() {
	appConfig := config.NewConfig()
	s.mainConfig = appConfig
	dbConnectionProvider := providers.NewDatabaseConnectionProvider(appConfig)
	repositoryProvider := providers.NewRepositoryProvider(dbConnectionProvider)
	stripeClientProvider := providers.NewClientStripeProvider(appConfig)
	stripeServiceProvider := providers.NewStripeServiceProvider(stripeClientProvider)
	serviceProvider := providers.NewServiceProvider(
		appConfig,
		repositoryProvider,
		dbConnectionProvider,
		stripeServiceProvider,
	)
	s.serviceProvider = serviceProvider
	grpcServerProvider := providers.NewGrpcServerProvider(serviceProvider)

	s.userGrpcServer = grpcServerProvider.UserGrpcServer()
	s.walletGrpcServer = grpcServerProvider.WalletGrpcServer()
	s.cardGrpcServer = grpcServerProvider.CardGrpcServer()

	s.grpcMiddlewareProvider = providers.NewGrpcMiddlewareProvider(serviceProvider)
}

func (s *Server) Start() error {
	s.initApp()
	server := grpc.NewServer(
		grpcMiddleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			grpcRecovery.UnaryServerInterceptor(),
			s.grpcMiddlewareProvider.Card().CreateCardAuthenticated,
			s.grpcMiddlewareProvider.Wallet().CreateWalletAuthenticated,
			s.grpcMiddlewareProvider.Wallet().WalletsListAuthenticated,
		),
	)
	user_server.RegisterUserServiceServer(server, s.userGrpcServer)
	walletServer.RegisterWalletServiceServer(server, s.walletGrpcServer)
	cardServer.RegisterCardServiceServer(server, s.cardGrpcServer)

	listener, e := net.Listen("tcp", fmt.Sprintf(":%s", s.mainConfig.GrpcPort))
	if e != nil {
		return e
	}
	return server.Serve(listener)
}
