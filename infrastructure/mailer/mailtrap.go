package mailer

import (
	"github.com/bccfilkom/drophere-go/domain"

	"gopkg.in/gomail.v2"
)

type mailtrap struct {
	dialer *gomail.Dialer
}

// NewMailtrap returns new mailtrap instance
func NewMailtrap(user, password string) domain.Mailer {
	dialer := gomail.NewDialer("smtp.mailtrap.io", 587, user, password)
	return &mailtrap{dialer}
}

// Send sends the email to mailtrap server
func (m *mailtrap) Send(from, to domain.MailAddress, subject, messagePlain, messageHTML string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", from.Address, from.Name)
	msg.SetAddressHeader("To", to.Address, to.Name)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", messageHTML)
	msg.AddAlternative("text/plain", messagePlain)

	return m.dialer.DialAndSend(msg)
}
