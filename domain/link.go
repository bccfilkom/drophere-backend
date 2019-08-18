package domain

import (
	"errors"
	"time"
)

var (
	// ErrLinkDuplicatedSlug error
	ErrLinkDuplicatedSlug = errors.New("Duplicated slug")
	// ErrLinkInvalidPassword error
	ErrLinkInvalidPassword = errors.New("Invalid password")
	// ErrLinkNotFound error
	ErrLinkNotFound = errors.New("Not found")
)

// Link domain model
type Link struct {
	ID                      uint
	UserID                  uint
	User                    *User
	Title                   string
	Password                string
	Slug                    string
	Deadline                *time.Time
	Description             string
	UserStorageCredentialID *uint
	UserStorageCredential   *UserStorageCredential
}

// IsProtected checks if the link is protected with password
func (l *Link) IsProtected() bool {
	return l.Password != ""
}

// LinkService abstraction
type LinkService interface {
	CheckLinkPassword(l *Link, password string) bool
	CreateLink(title, slug, description string, deadline *time.Time, password *string, user *User, providerID *uint) (*Link, error)
	UpdateLink(id uint, title, slug string, description *string, deadline *time.Time, password *string, providerID *uint) (*Link, error)
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
