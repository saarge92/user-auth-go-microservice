package app

import (
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/payment"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/providers"
	"go-user-microservice/pkg/protobuf/core"
	"go-user-microservice/scripts"
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
	paymentGrpcServer      *payment.GrpcServerPayment
	grpcMiddlewareProvider *providers.GrpcMiddlewareProvider
}

func NewServer() *Server {
	mainServer := &Server{}
	return mainServer
}
func (s *Server) InitConfig() error {
	var _, filename, _, _ = runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")

	if e := os.Chdir(dir); e != nil {
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
	scripts.Migrate(appConfig)
	dbConnectionProvider := providers.NewDatabaseConnectionProvider(appConfig)
	repositoryProvider := providers.NewRepositoryProvider(dbConnectionProvider)
	stripeClientProvider := providers.NewClientStripeProvider(appConfig)
	stripeServiceProvider := providers.NewStripeServiceProvider(stripeClientProvider)
	serviceProvider := providers.NewServiceProvider(
		appConfig,
		repositoryProvider,
		stripeServiceProvider,
	)
	transactionHandler := database.NewTransactionHandler(dbConnectionProvider.GetCoreConnection())

	grpcServerProvider := providers.NewGrpcServerProvider(serviceProvider, transactionHandler)

	s.userGrpcServer = grpcServerProvider.UserGrpcServer()
	s.walletGrpcServer = grpcServerProvider.WalletGrpcServer()
	s.cardGrpcServer = grpcServerProvider.CardGrpcServer()
	s.paymentGrpcServer = grpcServerProvider.PaymentGrpcServer()
	s.grpcMiddlewareProvider = providers.NewGrpcMiddlewareProvider(serviceProvider)
}

func (s *Server) Start() error {
	s.initApp()
	server := grpc.NewServer(
		grpcMiddleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			grpcRecovery.UnaryServerInterceptor(),
			s.grpcMiddlewareProvider.Card().CardsRequestAuthenticated,
			s.grpcMiddlewareProvider.Wallet().WalletsRequestsAuthenticated,
			s.grpcMiddlewareProvider.Payment().PaymentsRequestsAuthenticated,
		),
	)
	core.RegisterUserServiceServer(server, s.userGrpcServer)
	core.RegisterWalletServiceServer(server, s.walletGrpcServer)
	core.RegisterCardServiceServer(server, s.cardGrpcServer)
	core.RegisterPaymentServiceServer(server, s.paymentGrpcServer)

	listener, e := net.Listen("tcp", fmt.Sprintf(":%s", s.mainConfig.GrpcPort))
	if e != nil {
		return e
	}
	return server.Serve(listener)
}
