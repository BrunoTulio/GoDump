package alert_dump

import "github.com/BrunoTulio/GoDump/pkg/mail"

type (
	AlertCommand struct {
		Message string
		Mail    []string
	}

	AlertDumpService interface {
		Execute(command AlertCommand) error
	}

	alertDumpMailService struct {
		mail mail.Mail
	}
)

// Execute implements AlertDumpService.
func (a *alertDumpMailService) Execute(command AlertCommand) error {
	return a.mail.Send(mail.Message{
		Recipient: command.Mail,
		Subject:   "GoDump Backup",
		Body:      command.Message,
	})
}

func New(mail mail.Mail) AlertDumpService {
	return &alertDumpMailService{
		mail: mail,
	}
}
