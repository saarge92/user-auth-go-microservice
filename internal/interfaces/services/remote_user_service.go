package services

type RemoteUserServiceInterface interface {
	CheckRemoteUser(login string) bool
}
