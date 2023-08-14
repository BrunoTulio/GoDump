package factory

import (
	"log"

	"github.com/BrunoTulio/GoDump/internal/service/alert_dump"
	"github.com/BrunoTulio/GoDump/internal/service/backup_files"
	"github.com/BrunoTulio/GoDump/internal/service/configuration_file"
	"github.com/BrunoTulio/GoDump/internal/service/drive"
	"github.com/BrunoTulio/GoDump/internal/service/drive_token"

	"github.com/BrunoTulio/GoDump/internal/service/process_dump"
)

func MakeUploadDriver() {

}

func MakeAlterDumpService() alert_dump.AlertDumpService {
	config, err := MakeConfigurationFileService().Get()

	if err != nil {
		log.Fatal(err)
	}

	return alert_dump.New(MakeMail(config.Alert.Mail))
}

func MakeConfigurationFileService() configuration_file.ConfigurationFileService {
	return configuration_file.New()
}

func MakeProcessDumpService() process_dump.ProcessDumpService {
	return process_dump.New()
}

func MakeDriveTokenService() drive_token.DriveTokenService {
	config, err := MakeOAuth2Config()

	if err != nil {
		log.Fatalf("Config from file json. Err: %v\n", err)
	}

	return drive_token.New(config)
}

func MakeBackupsFilesService() backup_files.BackupFilesService {
	return backup_files.New()
}

func MakeDriveService() drive.DriveService {
	service, err := MakeGoogleDriveService()

	if err != nil {
		log.Fatal(err)
	}

	return drive.New(service)
}
