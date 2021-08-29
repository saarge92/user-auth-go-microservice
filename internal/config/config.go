package config

import "os"

type Config struct {
	GrpcPort          string
	CoreDatabaseUrl   string
	RemoteUserUrl     string
	AuthUserRemoteKey string
}

func NewConfig() *Config {
	return &Config{
		GrpcPort:          os.Getenv("GRPC_PORT"),
		CoreDatabaseUrl:   os.Getenv("CORE_DATABASE_URL"),
		RemoteUserUrl:     os.Getenv("REMOTE_USER_URL"),
		AuthUserRemoteKey: os.Getenv("AUTH_REMOTE_KEY"),
	}
}
