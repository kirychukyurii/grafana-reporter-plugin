package mail

import (
	"crypto/tls"
	"fmt"
	"net/mail"

	gomail "gopkg.in/mail.v2"
)

type Sender interface {
	Send(to []string, subject, body []byte, attachments []string) error
}

type Mail struct {
	from string

	Dialer *gomail.Dialer
}

func New(host string, port int, username, password string) (*Mail, error) {
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	dial, err := d.Dial()
	if err != nil {
		return nil, err
	}
	defer dial.Close()

	from := username
	_, err = mail.ParseAddress(from)
	if err != nil {
		from = fmt.Sprintf("%s@%s", username, host)
	}

	return &Mail{
		from: from,

		Dialer: d,
	}, nil
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
