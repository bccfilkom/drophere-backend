package hasher

import (
	"github.com/bccfilkom/drophere-go/domain"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

// NewBcryptHasher bcrypt hasher that implements Hasher interface
func NewBcryptHasher() domain.Hasher {
	return &bcryptHasher{
		cost: 12,
	}
}

// Hash implementation
func (b *bcryptHasher) Hash(s string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), b.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify implementation
func (b *bcryptHasher) Verify(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
