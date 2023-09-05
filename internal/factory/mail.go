package factory

import (
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/pkg/mail"
	"github.com/BrunoTulio/GoDump/pkg/mail/gomail"
	"github.com/BrunoTulio/GoDump/pkg/mail/smtp"
)

func MakeSMTPMail(mail domain.Mail) mail.Mail {
	return smtp.NewWithOption(
		smtp.NewOptionWithParmas(
			mail.User,
			mail.Password,
			mail.Host,
			mail.Port,
			mail.InsureSecurity,
			mail.From,
		),
	)
}

func MakeGoMail(mail domain.Mail) mail.Mail {
	return gomail.NewWithOption(
		gomail.NewOptionWithParmas(
			mail.User,
			mail.Password,
			mail.Host,
			mail.Port,
			mail.InsureSecurity,
			mail.From,
		),
	)
}
