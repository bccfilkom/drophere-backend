package domain

import "errors"

// ErrTemplateNotFound error
var ErrTemplateNotFound = errors.New("Template not found")

// MailAddress model
type MailAddress struct {
	Address string
	Name    string
}

// Mailer abstraction
type Mailer interface {
	Send(from, to MailAddress, subject, messagePlain, messageHTML string) error
}
