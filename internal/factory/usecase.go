package factory

import (
	"github.com/BrunoTulio/GoDump/internal/application/usecase/backup"
	"github.com/BrunoTulio/GoDump/internal/application/usecase/restore"
	"github.com/BrunoTulio/GoDump/internal/service/alert_dump"
	"github.com/BrunoTulio/GoDump/internal/service/backup_files"
	"github.com/BrunoTulio/GoDump/internal/service/configuration_file"
	"github.com/BrunoTulio/GoDump/internal/service/drive"
	"github.com/BrunoTulio/GoDump/internal/service/process_dump"

	"github.com/robfig/cron/v3"
)

func MakeBackupUseCase(
	configurationFileService configuration_file.ConfigurationFileService,
	processDumpService process_dump.ProcessDumpService,

) backup.BackupUseCase {
	return backup.NewBackup(configurationFileService, processDumpService)
}

func MakeDefaultBackupUseCase() backup.BackupUseCase {
	return backup.NewBackup(
		MakeConfigurationFileService(),
		MakeProcessDumpService(),
	)
}

func MakeBackupAutoUseCase(
	configurationFileService configuration_file.ConfigurationFileService,
	processDumpService process_dump.ProcessDumpService,
	alertDumpService alert_dump.AlertDumpService,
	driveService drive.DriveService,
	backupsFilesService backup_files.BackupFilesService,
	cron *cron.Cron,
) backup.BackupAutoUseCase {
	return backup.NewBackupAuto(
		configurationFileService,
		processDumpService,
		alertDumpService,
		driveService,
		backupsFilesService,
		cron,
	)
}

func MakeDefaultBackupAutoUseCase(cron *cron.Cron) backup.BackupAutoUseCase {
	return backup.NewBackupAuto(
		MakeConfigurationFileService(),
		MakeProcessDumpService(),
		MakeAlterDumpService(),
		MakeDriveService(),
		MakeBackupsFilesService(),
		cron,
	)
}

func MakeRestoreAutoUseCase(
	configurationFileService configuration_file.ConfigurationFileService,
	backupFilesService backup_files.BackupFilesService,
	processDumpService process_dump.ProcessDumpService,
	driveService drive.DriveService,
) restore.RestoreAutoUseCase {
	return restore.NewRestoreAuto(
		configurationFileService,
		backupFilesService,
		processDumpService,
		driveService,
	)
}

func MakeDefaultRestoreAutoUseCase() restore.RestoreAutoUseCase {
	return restore.NewRestoreAuto(
		MakeConfigurationFileService(),
		MakeBackupsFilesService(),
		MakeProcessDumpService(),
		MakeDriveService(),
	)
}
