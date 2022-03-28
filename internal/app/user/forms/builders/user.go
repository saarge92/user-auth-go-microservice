package builders

import (
	"go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/pkg/protobuf/core"
)

type UserFormBuilder struct {
	userRepository domain.UserRepository
}

func NewUserFormBuilder(userRepository domain.UserRepository) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *core.SignUpMessage) *forms.SignUp {
	return forms.NewSignUpForm(request)
}

func (b *UserFormBuilder) SignIn(request *core.SignInMessage) *forms.SignIn {
	return &forms.SignIn{SignInMessage: request}
}
