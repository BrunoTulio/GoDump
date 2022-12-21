package listener

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/internal/settings"
	"github.com/BrunoTulio/GoDump/pkg/event"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/BrunoTulio/GoDump/pkg/mail"
	"github.com/pkg/errors"
)

type BackupListener interface {
	event.Listener
}

type backupListener struct {
	backupCheckDumb       backup.BackupCheckDump
	backupCheckConnection backup.BackupCheckConnection
	backupMake            backup.BackupMake
	backupRepository      backup.BackupRepository
	mail                  mail.Mail
}

func (d backupListener) Handler(ctx context.Context, e event.Event) error {

	switch e.(type) {
	case *backup.BackupLocalEvent:
		{
			_, err := d.makeBackup(ctx, settings.DatabaseName())

			if err != nil {
				return err
			}

			err = d.backupRepository.StoreLastFile(time.Now())

			if err != nil {
				logger.Warn(err)
			}

			return nil
		}
	case *backup.BackupMailEvent:
		{
			path, err := d.makeBackup(ctx, fmt.Sprintf("%s.mail", settings.DatabaseName()))

			if err != nil {
				return err
			}

			bytes, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			message := mail.Message{
				Recipient: []string{settings.MailSend()},
				From:      "godump@suport.com.br",
				Subject:   "Backup successfully concluded",
				Body:      `<p>GoDump dump finish</p>`,
				Attachments: map[string][]byte{
					settings.DatabaseName(): bytes,
				},
			}

			err = d.mail.Send(message)

			if err != nil {
				return err
			}

			err = d.backupRepository.StoreLastMail(time.Now())

			if err != nil {
				logger.Warn(err)
			}

			return nil
		}
	default:
		return nil
	}

}

func (d backupListener) makeBackup(ctx context.Context, nameFile string) (string, error) {
	err := d.backupCheckDumb.Execute()

	if err != nil {
		return "", errors.Wrap(err, "client postgres dump not found")
	}

	err = d.backupCheckConnection.Execute(ctx)

	if err != nil {
		return "", errors.Wrap(err, "connection fail")
	}

	path, err := d.backupMake.Execute(ctx, nameFile)

	if err != nil {
		return "", errors.Wrap(err, "make backup fail")
	}

	return path, nil
}

func NewBackupListener(
	backupCheckDumb backup.BackupCheckDump,
	backupCheckConnection backup.BackupCheckConnection,
	backupMake backup.BackupMake,
	backupRepository backup.BackupRepository,
	mail mail.Mail,
) BackupListener {
	return backupListener{
		backupCheckDumb:       backupCheckDumb,
		backupCheckConnection: backupCheckConnection,
		backupMake:            backupMake,
		mail:                  mail,
		backupRepository:      backupRepository,
	}
}
