package services

type UserRemoteMock struct {
}

func (m *UserRemoteMock) CheckRemoteUser(inn uint64) (r bool, e error) {
	return true, nil
}
