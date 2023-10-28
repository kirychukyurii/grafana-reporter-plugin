package smtp

import (
	"crypto/tls"
	"fmt"
	"net/mail"

	gomail "gopkg.in/mail.v2"
)

type Mail struct {
	from string

	Dialer *gomail.Dialer
}

type Options struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func New(options *Options) (*Mail, error) {
	d := gomail.NewDialer(options.Host, options.Port, options.Username, options.Password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	/*dial, err := d.Dial()
	if err != nil {
		return nil, err
	}
	defer dial.Close()*/

	from := options.Username
	_, err := mail.ParseAddress(from)
	if err != nil {
		from = fmt.Sprintf("%s@%s", options.Username, options.Host)
	}

	return &Mail{
		from: from,

		Dialer: d,
	}, nil
}
