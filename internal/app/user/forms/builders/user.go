package builders

import (
	"go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/pkg/protobuf/user_server"
)

type UserFormBuilder struct {
	userRepository domain.UserRepository
}

func NewUserFormBuilder(userRepository domain.UserRepository) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *user_server.SignUpMessage) *forms.SignUp {
	return forms.NewSignUpForm(request)
}

func (b *UserFormBuilder) SignIn(request *user_server.SignInMessage) *forms.SignIn {
	return &forms.SignIn{SignInMessage: request}
}
