package config

import "os"

type Config struct {
	GrpcPort          string
	CoreDatabaseURL   string
	RemoteUserURL     string
	AuthUserRemoteKey string
}

func NewConfig() *Config {
	return &Config{
		GrpcPort:          os.Getenv("GRPC_PORT"),
		CoreDatabaseURL:   os.Getenv("CORE_DATABASE_URL"),
		RemoteUserURL:     os.Getenv("REMOTE_USER_URL"),
		AuthUserRemoteKey: os.Getenv("AUTH_REMOTE_KEY"),
	}
}
