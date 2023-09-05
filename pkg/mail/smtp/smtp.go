package smtp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/BrunoTulio/GoDump/pkg/mail"
)

type mailSmtp struct {
	*Options
}

// Send implements mail.Mail.
func (m mailSmtp) Send(message mail.Message) error {
	data, err := m.data(message)

	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.User, m.Password, m.Host)

	if m.InsureSecurity {
		err = m.sendInsecureSkip(auth, message, data)
		return err
	}

	return m.send(auth, message, data)
}

func (e mailSmtp) send(auth smtp.Auth, message mail.Message, data []byte) error {
	return smtp.SendMail(fmt.Sprintf("%s:%d", e.Host, e.Port), auth, e.User,
		message.Recipient, data)
}

func (e mailSmtp) sendInsecureSkip(auth smtp.Auth, message mail.Message, data []byte) (err error) {
	smtpConn, err := smtp.Dial(fmt.Sprintf("%s:%d", e.Host, e.Port))

	if err != nil {
		return
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	err = smtpConn.StartTLS(tlsConfig)

	if err != nil {
		return
	}

	defer func() {
		err = smtpConn.Quit()
		logger.Error(err)
	}()

	if err = smtpConn.Auth(auth); err != nil {
		return
	}

	// To && From
	if err = smtpConn.Mail(e.From); err != nil {
		return
	}

	recs := strings.Join(message.Recipient, ",")
	if err = smtpConn.Rcpt(recs); err != nil {
		return
	}

	// Data
	w, err := smtpConn.Data()
	if err != nil {
		return
	}

	defer func() {
		err = w.Close()
		logger.Error(err)
	}()

	_, err = w.Write(data)
	if err != nil {
		return
	}

	return
}

func (e mailSmtp) data(message mail.Message) ([]byte, error) {

	withAttachments := len(message.Attachments) > 0

	data := message.Body
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("Subject: %s\n", message.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(message.Recipient, ",")))

	buf.WriteString(mimeHeaders)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	buf.WriteString(data)
	if withAttachments {
		for k, v := range message.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
		buf.WriteString("--")
	}
	return buf.Bytes(), nil

}

func NewWithOption(o *Options) mail.Mail {
	return mailSmtp{Options: o}
}

func New() mail.Mail {
	return NewWithOption(NewOption())
}
