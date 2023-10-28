package smtp

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

type Sender interface {
	Send(to []string, subject, body []byte, attachments []string) error
}

func (m *Mail) Send(to []string, subject, body []byte, attachments []string) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", string(subject))
	msg.SetBody("text/html", string(body))

	if len(attachments) > 0 {
		for _, a := range attachments {
			msg.Attach(a)
		}
	}

	if err := m.Dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("send message: %v", err)
	}

	return nil
}
