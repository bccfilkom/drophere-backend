package hasher

import "github.com/bccfilkom/drophere-go/domain"

type notAHasher struct{}

// NewNotAHasher returns notAHasher instance for testing purpose
func NewNotAHasher() domain.Hasher {
	return &notAHasher{}
}

// Hash implementation
func (n *notAHasher) Hash(s string) (string, error) {
	return s, nil
}

// Verify implementation
func (n *notAHasher) Verify(hashed, plain string) bool {
	return hashed == plain
}
