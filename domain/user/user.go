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
func (s *service) Register(email, name, password string) (*domain.User, error) {
	user := &domain.User{
		Email: email,
		Name:  name,
	}
	user.SetPassword(password)
	return s.userRepo.Create(user)
}

// Auth implementation
func (s *service) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password) {
		return nil, domain.ErrUserInvalidPassword
	}

	return s.authenticator.Authenticate(user)
}

// Update implementation
func (s *service) Update(userID uint, name, newPassword, oldPassword *string) (*domain.User, error) {
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if newPassword != nil {
		if oldPassword == nil || !u.VerifyPassword(*oldPassword) {
			return nil, domain.ErrUserInvalidPassword
		}

		u.SetPassword(*newPassword)
	}

	if name != nil {
		u.Name = *name
	}

	return s.userRepo.Update(u)
}
