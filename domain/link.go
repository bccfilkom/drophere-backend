package domain

import (
	"errors"
	"time"
)

var (
	// ErrLinkDuplicatedSlug error
	ErrLinkDuplicatedSlug = errors.New("link: duplicated slug")
	// ErrLinkNotFound error
	ErrLinkNotFound = errors.New("link: not found")
)

// Link domain model
type Link struct {
	ID          uint
	UserID      uint
	User        *User
	Title       string
	Password    string
	Slug        string
	Deadline    time.Time
	Description string
}

// IsProtected checks if the link is protected with password
func (l *Link) IsProtected() bool {
	return l.Password != ""
}

// SetPassword hash input password and set it to the link struct
func (l *Link) SetPassword(password string) {
	l.Password = password
}

// VerifyPassword checks if the encrypted password content is
// equal to the given plain password
func (l *Link) VerifyPassword(plainPwd string) bool {
	return l.Password == plainPwd
}

// LinkService abstraction
type LinkService interface {
	CreateLink(title, slug, description string, user *User) (*Link, error)
}

// LinkRepository abstraction
type LinkRepository interface {
	Create(l *Link) (*Link, error)
	FindBySlug(slug string) (*Link, error)
}
