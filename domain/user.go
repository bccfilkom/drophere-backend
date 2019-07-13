package domain

import (
	"errors"
	"time"
)

var (
	// ErrUserDuplicated error
	ErrUserDuplicated = errors.New("Duplicated email")
	// ErrUserInvalidPassword error
	ErrUserInvalidPassword = errors.New("Invalid password")
	// ErrUserNotFound error
	ErrUserNotFound = errors.New("User not found")
)

// User model
type User struct {
	ID           uint
	Email        string
	Name         string
	Password     string
	DropboxToken *string
	DriveToken   *string
}

// UserCredentials model
type UserCredentials struct {
	Token  string
	Expiry *time.Time
}

// UserService abstraction
type UserService interface {
	Register(email, name, password string) (*User, error)
	Auth(email, password string) (*UserCredentials, error)
	Update(userID uint, name, password, oldPassword *string) (*User, error)
	UpdateStorageToken(userID uint, dropboxToken *string) (*User, error)
}

// UserRepository abstraction
type UserRepository interface {
	Create(u *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(u *User) (*User, error)
}

// Authenticator is external authentication service
type Authenticator interface {
	Authenticate(u *User) (*UserCredentials, error)
}
