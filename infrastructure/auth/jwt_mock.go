package auth

import (
	"strconv"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

type jwtAuthenticatorMock struct{}

// NewJWTMock func
func NewJWTMock() domain.Authenticator {
	return &jwtAuthenticatorMock{}
}

// Authenticate mock
func (j *jwtAuthenticatorMock) Authenticate(u *domain.User) (*domain.UserCredentials, error) {
	t := time.Now().Add(time.Hour)
	return &domain.UserCredentials{
		Token: "user_token_"+strconv.Itoa(int(u.ID)),
		Expiry: &t,
	}, nil
}