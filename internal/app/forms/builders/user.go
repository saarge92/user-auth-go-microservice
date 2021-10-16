package builders

import (
	repositories2 "go-user-microservice/internal/app/domain/repositories"
	forms2 "go-user-microservice/internal/app/forms"
	functions2 "go-user-microservice/internal/app/forms/functions"
	"go-user-microservice/pkg/protobuf/user"
)

type UserFormBuilder struct {
	userRepository repositories2.UserRepositoryInterface
}

func NewUserFormBuilder(userRepository repositories2.UserRepositoryInterface) *UserFormBuilder {
	return &UserFormBuilder{
		userRepository: userRepository,
	}
}

func (b *UserFormBuilder) Signup(request *user.SignUpMessage) *forms2.SignUp {
	userExistRule := functions2.UserAlReadyExists(b.userRepository)
	userInnRule := functions2.UserWithInnAlreadyExists(b.userRepository)
	return forms2.NewSignUpForm(request, userExistRule, userInnRule)
}

func (b *UserFormBuilder) SignIn(request *user.SignInMessage) *forms2.SignIn {
	return &forms2.SignIn{SignInMessage: request}
}
