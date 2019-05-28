package domain

import (
	"time"
)

// Link domain model
type Link struct {
	ID          uint
	UserID      uint
	User        *User
	Title       string
	Password    *string
	Slug        string
	Deadline    time.Time
	Description string
}

// IsProtected checks if the link is protected with password
func (l *Link) IsProtected() bool {
	return l.Password != nil
}
