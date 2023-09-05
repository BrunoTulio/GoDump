package gomail

import (
	"crypto/tls"

	"github.com/BrunoTulio/GoDump/pkg/mail"
	goMail "gopkg.in/mail.v2"
)

type goMailSmtp struct {
	*Options
}

// Send implements mail.Mail.
func (g goMailSmtp) Send(message mail.Message) error {
	d := goMail.NewDialer(g.Host, g.Port, g.User, g.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: g.InsureSecurity}

	m := goMail.NewMessage()
	m.SetHeader("From", g.From)
	m.SetHeader("To", message.Recipient...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/plain", message.Body)

	return d.DialAndSend(m)
}

func NewWithOption(o *Options) mail.Mail {
	return goMailSmtp{Options: o}
}

func New() mail.Mail {
	return NewWithOption(NewOption())
}
