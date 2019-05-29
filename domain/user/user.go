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

	var err error
	user.Password, err = s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return s.userRepo.Create(user)
}

// Auth implementation
func (s *service) Auth(email, password string) (*domain.UserCredentials, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !s.passwordHasher.Verify(user.Password, password) {
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
		if oldPassword == nil || !s.passwordHasher.Verify(u.Password, *oldPassword) {
			return nil, domain.ErrUserInvalidPassword
		}

		u.Password, err = s.passwordHasher.Hash(*newPassword)
		if err != nil {
			return nil, err
		}
	}

	if name != nil {
		u.Name = *name
	}

	return s.userRepo.Update(u)
}
