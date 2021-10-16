package builders

import (
	"go-user-microservice/internal/app/domain/repositories"
	functions2 "go-user-microservice/internal/app/forms/functions"
	formUser "go-user-microservice/internal/app/forms/user"
	"go-user-microservice/pkg/protobuf/user"
)

type UserFormBuilder struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserFormBuilder(userRepository repositories.UserRepositoryInterface) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *user.SignUpMessage) *formUser.SignUp {
	userExistRule := functions2.UserAlReadyExists(b.userRepository)
	userInnRule := functions2.UserWithInnAlreadyExists(b.userRepository)
	return formUser.NewSignUpForm(request, userExistRule, userInnRule)
}

func (b *UserFormBuilder) SignIn(request *user.SignInMessage) *formUser.SignIn {
	return &formUser.SignIn{SignInMessage: request}
}
