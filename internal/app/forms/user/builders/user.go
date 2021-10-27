package builders

import (
	"go-user-microservice/internal/app/domain/repositories"
	formUser "go-user-microservice/internal/app/forms/user"
	"go-user-microservice/pkg/protobuf/user_server"
)

type UserFormBuilder struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserFormBuilder(userRepository repositories.UserRepositoryInterface) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *user_server.SignUpMessage) *formUser.SignUp {
	return formUser.NewSignUpForm(request)
}

func (b *UserFormBuilder) SignIn(request *user_server.SignInMessage) *formUser.SignIn {
	return &formUser.SignIn{SignInMessage: request}
}
