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
	Deadline    *time.Time
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
	UpdateLink(id uint, title, slug string, description *string, deadline *time.Time, password *string) (*Link, error)
	DeleteLink(id uint) error
	FetchLink(id uint) (*Link, error)
	FindLinkBySlug(slug string) (*Link, error)
	ListLinks(userID uint) ([]Link, error)
}

// LinkRepository abstraction
type LinkRepository interface {
	Create(l *Link) (*Link, error)
	Delete(l *Link) error
	FindByID(id uint) (*Link, error)
	FindBySlug(slug string) (*Link, error)
	ListByUser(userID uint) ([]Link, error)
	Update(l *Link) (*Link, error)
}
