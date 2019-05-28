package user

import (
	"github.com/bccfilkom/drophere-go/domain"
)

type service struct {
	userRepo      domain.UserRepository
	authenticator domain.Authenticator
}

// NewService returns service instance
func NewService(userRepo domain.UserRepository, authenticator domain.Authenticator) domain.UserService {
	return &service{
		userRepo:      userRepo,
		authenticator: authenticator,
	}
}

// Register implementation
func (u *service) Register(email, name, password string) (*domain.User, error) {
	user := &domain.User{
		Email: email,
		Name:  name,
	}
	user.SetPassword(password)
	return u.userRepo.Create(user)
}

// Auth implementation
func (u *service) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, domain.ErrUserInvalidPassword
	}

	return u.authenticator.Authenticate(user)
}
