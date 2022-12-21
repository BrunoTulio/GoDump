package simplemail

import (
	"crypto/tls"
	"fmt"
	"time"

	myMail "github.com/BrunoTulio/GoDump/pkg/mail"
	mail "github.com/xhit/go-simple-mail/v2"
)

type simpleMail struct {
	host     string
	user     string
	password string
	port     int
	secure   bool
}

// Send implements mail.Mail
func (s simpleMail) Send(message myMail.Message) error {
	server := mail.NewSMTPClient()
	server.Host = s.host
	server.Port = s.port
	server.Username = s.user
	server.Password = s.password
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("Go Dumb <%s>", message.From)).
		AddTo(message.Recipient...).
		//AddCc(message.Recipient...).
		SetSubject(message.Subject)

	email.SetBody(mail.TextHTML, message.Body)

	for name, value := range message.Attachments {
		email.Attach(&mail.File{Data: value, Name: name})
	}
	if email.Error != nil {
		return email.Error
	}

	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil

}

func New(
	host string,
	user string,
	password string,
	port int,
	secure bool,
) myMail.Mail {
	return simpleMail{
		host:     host,
		user:     user,
		password: password,
		port:     port,
		secure:   secure,
	}
}
