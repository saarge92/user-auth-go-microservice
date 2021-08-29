package services

type RemoteUserService struct {
}

func NewRemoteUserService() *RemoteUserService {
	return &RemoteUserService{}
}

func (s *RemoteUserService) CheckRemoteUser(login string) bool {
	return true
}
