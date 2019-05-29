package domain

// Hasher abstraction
type Hasher interface {
	Hash(s string) (string, error)
	Verify(hashed, plain string) bool
}