package domain

import "errors"

var (
	// ErrUserStorageCredentialNotFound error
	ErrUserStorageCredentialNotFound = errors.New("User Storage Credential not found")
)

// UserStorageCredential stores information about user's account on
// a storage provider (e.g. Dropbox)
type UserStorageCredential struct {
	ID                 uint
	UserID             uint
	User               User
	ProviderID         uint
	ProviderCredential string
	Email              string
	Photo              string
}

// UserStorageCredentialFilters stores filters to be used by
// Find function in UserStorageCredentialRepository
type UserStorageCredentialFilters struct {
	UserIDs     []uint
	ProviderIDs []uint
}

// UserStorageCredentialRepository abstraction
type UserStorageCredentialRepository interface {
	Find(filters UserStorageCredentialFilters, withUserRelation bool) ([]UserStorageCredential, error)
	FindByID(id uint, withUserRelation bool) (UserStorageCredential, error)
	Create(cred UserStorageCredential) (UserStorageCredential, error)
	Update(cred UserStorageCredential) (UserStorageCredential, error)
	Delete(cred UserStorageCredential) error
}
