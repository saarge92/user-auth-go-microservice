package servers

type ServerInterface interface {
	InitConfig() error
	InitContainer() error
	Start() error
	InitApp() error
}
