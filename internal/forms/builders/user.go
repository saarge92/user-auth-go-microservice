package builders

import (
	"go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/forms"
	"go-user-microservice/internal/forms/functions"
	"go-user-microservice/pkg/protobuf/user"
)

type UserFormBuilder struct {
	userRepository repositories.UserRepository
}

func NewUserFormBuilder(userRepository repositories.UserRepository) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *user.SignUpMessage) *forms.SignUp {
	userExistRule := functions.UserAlReadyExists(b.userRepository)
	return forms.NewSignUpForm(request, userExistRule)
}

func (b *UserFormBuilder) SignIn(request *user.SignInMessage) *forms.SignIn {
	return &forms.SignIn{SignInMessage: request}
}
