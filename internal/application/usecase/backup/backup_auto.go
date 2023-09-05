package backup

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/internal/service/alert_dump"
	"github.com/BrunoTulio/GoDump/internal/service/backup_files"
	"github.com/BrunoTulio/GoDump/internal/service/configuration_file"
	"github.com/BrunoTulio/GoDump/internal/service/drive"

	"github.com/BrunoTulio/GoDump/internal/service/process_dump"

	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

type (
	BackupAutoUseCase interface {
		Generate(ctx context.Context)
	}

	backupAutoUseCase struct {
		configurationFileService configuration_file.ConfigurationFileService
		processDumpService       process_dump.ProcessDumpService
		alertDumpService         alert_dump.AlertDumpService
		driveService             drive.DriveService
		backupsFilesService      backup_files.BackupFilesService
		cron                     *cron.Cron
	}
)

// Execute implements BackupGenerateUseCase.
func (b *backupAutoUseCase) Generate(ctx context.Context) {

	config, err := b.configurationFileService.Get()

	if err != nil {
		logger.Fatal(err)
	}

	err = config.IsValid()

	if err != nil {
		logger.Fatal(err)
	}
	sort.Sort(config.Dumps)

	b.cron.AddFunc(config.Cron, func() {
		logger.Info("Initialize event backup database file")

		group := errgroup.Group{}
		for _, dump := range config.Dumps {
			dump := dump

			group.Go(func() error {

				b.backupsFilesService.RemoveByDuration(config.DurationFileInFolder, dump)

				pathFile, err := b.processDumpService.Backup(dump)

				if err != nil {
					b.notifyError(config, err)
					return err
				}

				b.notifySuccess(config, fmt.Sprintf("Sucesso backup dump service path file %s", pathFile))

				if !config.GoogleDrive.Enable {
					return nil
				}

				err = b.driveService.Upload(context.Background(), pathFile, dump)

				if err != nil {
					b.notifyError(config, err)
					return err
				}

				b.notifySuccess(config, "Sucesso backup dump service google drive")
				return nil
			})

			err := group.Wait()

			if err != nil {
				logger.Error(err)
			}
		}
	})

}

func (b *backupAutoUseCase) notifyError(config domain.Config, err error) {
	if !config.Alert.Enable {
		return
	}

	b.alertDumpService.Execute(alert_dump.AlertCommand{
		Message: fmt.Sprintf("Error GoDump %v %s", err, time.Now().Format(constants.LayoutDate)),
		Mail:    config.Alert.SendMail,
	})
}

func (b *backupAutoUseCase) notifySuccess(config domain.Config, msg string) {
	if !config.Alert.Enable {
		return
	}

	b.alertDumpService.Execute(alert_dump.AlertCommand{
		Message: msg,
		Mail:    config.Alert.SendMail,
	})
}

func NewBackupAuto(
	configurationFileService configuration_file.ConfigurationFileService,
	processDumpService process_dump.ProcessDumpService,
	alertDumpService alert_dump.AlertDumpService,
	driveService drive.DriveService,
	backupsFilesService backup_files.BackupFilesService,
	cron *cron.Cron,
) BackupAutoUseCase {
	return &backupAutoUseCase{
		configurationFileService,
		processDumpService,
		alertDumpService,
		driveService,
		backupsFilesService,
		cron,
	}
}
