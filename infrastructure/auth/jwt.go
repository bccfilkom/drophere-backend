package auth

import (
	"time"

	"github.com/bccfilkom/drophere-go/domain"
	jwt "github.com/dgrijalva/jwt-go"
)

type jwtAuthenticator struct {
	key []byte

	duration time.Duration

	algo string
}

// NewJWT func
func NewJWT(secret string, duration time.Duration, algo string) domain.Authenticator {
	return &jwtAuthenticator{
		key:      []byte(secret),
		duration: duration,
		algo:     algo,
	}
}

// Authenticate func
func (j *jwtAuthenticator) Authenticate(u *domain.User) (*domain.UserCredentials, error) {
	expiry := time.Now().Add(j.duration)
	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.algo), jwt.MapClaims{
		"user_id": u.ID,
		"exp":     expiry.Unix(),
	})

	tokenS, err := token.SignedString(j.key)
	if err != nil {
		return nil, err
	}
	return &domain.UserCredentials{
		Token:  tokenS,
		Expiry: &expiry,
	}, nil
}
