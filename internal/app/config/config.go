package config

import (
	"os"
	"strconv"
)

type Config struct {
	GrpcPort            string
	CoreDatabaseURL     string
	RemoteUserURL       string
	AuthUserRemoteKey   string
	JwtAudience         string
	JwtExpiration       int32
	JwtKey              string
	DatabaseDriver      string
	SecretStripeKey     string
	SecretEncryptionKey string
}

func NewConfig() *Config {
	jwtExpiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	return &Config{
		GrpcPort:            os.Getenv("GRPC_PORT"),
		CoreDatabaseURL:     os.Getenv("CORE_DATABASE_URL"),
		RemoteUserURL:       os.Getenv("REMOTE_USER_URL"),
		AuthUserRemoteKey:   os.Getenv("AUTH_REMOTE_KEY"),
		JwtAudience:         os.Getenv("JWT_AUDIENCE"),
		JwtExpiration:       int32(jwtExpiration),
		JwtKey:              os.Getenv("JWT_KEY"),
		DatabaseDriver:      os.Getenv("DATABASE_DRIVER"),
		SecretStripeKey:     os.Getenv("STRIPE_SECRET_KEY"),
		SecretEncryptionKey: os.Getenv("ENCRYPTION_SECRET_KEY"),
	}
}
