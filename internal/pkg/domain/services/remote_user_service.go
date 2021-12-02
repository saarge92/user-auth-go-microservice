package services

type RemoteUserService interface {
	CheckRemoteUser(inn uint64) (r bool, e error)
}
