package mailer

import "github.com/bccfilkom/drophere-go/domain"

// MockMessage is for testing purpose
type MockMessage struct {
	From         string
	To           string
	Title        string
	MessagePlain string
	MessageHTML  string
}

// MockMessages is an in-memory storage for testing purpose
var MockMessages []MockMessage

func init() {
	MockMessages = make([]MockMessage, 0)
}

// ClearMessages reset the MockMessages
func ClearMessages() {
	MockMessages = make([]MockMessage, 0)
}

type mockMailer struct{}

// NewMockMailer returns new mockMailer instance
func NewMockMailer() domain.Mailer {
	return &mockMailer{}
}

// Send sends the email to mockMailer server
func (m *mockMailer) Send(from, to domain.MailAddress, subject, messagePlain, messageHTML string) error {
	MockMessages = append(MockMessages, MockMessage{from.Address, to.Address, subject, messagePlain, messageHTML})
	return nil
}
