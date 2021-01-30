package mailer

import (
	"log"

	"github.com/bccfilkom/drophere-go/domain"

	sendgridClient "github.com/sendgrid/sendgrid-go"
	mailHelper "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgrid struct {
	apiKey string
	debug  bool
}

// NewSendgrid returns new sendgrid instance
func NewSendgrid(apiKey string, debug bool) domain.Mailer {
	return &sendgrid{apiKey, debug}
}

// Send sends the email to sendgrid server
func (s *sendgrid) Send(from, to domain.MailAddress, subject, messagePlain, messageHTML string) error {

	f := mailHelper.NewEmail(
		from.Name,
		from.Address,
	)
	t := mailHelper.NewEmail(
		to.Name,
		to.Address,
	)
	msg := mailHelper.NewSingleEmail(f, subject, t, messagePlain, messageHTML)
	c := sendgridClient.NewSendClient(s.apiKey)

	// log the response
	resp, err := c.Send(msg)
	if s.debug {
		log.Println("Sendgrid: StatusCode:", resp.StatusCode)
		log.Println("Sendgrid: Body:", resp.Body)
		log.Println("Sendgrid: Headers:", resp.Headers)
	}

	return err

}
