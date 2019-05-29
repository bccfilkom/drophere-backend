package user

import (
	"github.com/bccfilkom/drophere-go/domain"
)

type service struct {
	userRepo       domain.UserRepository
	authenticator  domain.Authenticator
	passwordHasher domain.Hasher
}

// NewService returns service instance
func NewService(
	userRepo domain.UserRepository,
	authenticator domain.Authenticator,
	passwordHasher domain.Hasher,
) domain.UserService {
	return &service{
		userRepo:       userRepo,
		authenticator:  authenticator,
		passwordHasher: passwordHasher,
	}
}

// Register implementation
func (s *service) Register(email, name, password string) (*domain.User, error) {
	user := &domain.User{
		Email: email,
		Name:  name,
	}
	user.SetPassword(password, s.passwordHasher)
	return s.userRepo.Create(user)
}

// Auth implementation
func (s *service) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(password, s.passwordHasher) {
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
		if oldPassword == nil || !u.VerifyPassword(*oldPassword, s.passwordHasher) {
			return nil, domain.ErrUserInvalidPassword
		}

		u.SetPassword(*newPassword, s.passwordHasher)
	}

	if name != nil {
		u.Name = *name
	}

	return s.userRepo.Update(u)
}
