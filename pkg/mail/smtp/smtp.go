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
	host     string
	user     string
	password string
	port     int
	secure   bool
}

func (m mailSmtp) Send(message mail.Message) error {

	data, err := m.data(message)

	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.user, m.password, m.host)

	if m.secure {
		return m.sendInsecureSkip(auth, message, data)
	}

	return m.send(auth, message, data)
}

func (e mailSmtp) sendInsecureSkip(auth smtp.Auth, message mail.Message, data []byte) error {
	smtpConn, err := smtp.Dial(fmt.Sprintf("%s:%d", e.host, e.port))

	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	err = smtpConn.StartTLS(tlsConfig)

	if err != nil {
		return err
	}

	defer func() {
		err := smtpConn.Quit()
		if err != nil {
			logger.Error(err)
		}
	}()

	if err = smtpConn.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = smtpConn.Mail(message.From); err != nil {
		return err
	}

	recs := strings.Join(message.Recipient, ",")
	if err = smtpConn.Rcpt(recs); err != nil {
		return err
	}

	// Data
	w, err := smtpConn.Data()
	if err != nil {
		return err
	}

	defer func() {
		err := w.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (e mailSmtp) send(auth smtp.Auth, message mail.Message, data []byte) error {
	return smtp.SendMail(fmt.Sprintf("%s:%d", e.host, e.port), auth, e.user,
		message.Recipient, data)
}

func (e mailSmtp) data(message mail.Message) ([]byte, error) {

	withAttachments := len(message.Attachments) > 0

	data := message.Body
	mimeHeaders := "MIME-Version: 1.0\n"

	if !withAttachments {

		data = message.Body
		mimeHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("Subject: %s\n", message.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(message.Recipient, ",")))

	buf.WriteString(mimeHeaders)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		//
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	buf.WriteString(data)
	if withAttachments {
		for k, v := range message.Attachments {
			contentType := http.DetectContentType(v)

			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", contentType))
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

func New(host, user, password string, port int, secure bool) mail.Mail {
	return mailSmtp{
		host:     host,
		user:     user,
		password: password,
		port:     port,
		secure:   secure,
	}
}
