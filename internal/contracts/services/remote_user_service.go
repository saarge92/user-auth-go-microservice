package services

type RemoteUserServiceInterface interface {
	CheckRemoteUser(inn uint64) (r bool, e error)
}
