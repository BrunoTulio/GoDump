package restore

import (
	"context"
	"os"
	"sort"

	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/internal/service/backup_files"
	"github.com/BrunoTulio/GoDump/internal/service/configuration_file"
	"github.com/BrunoTulio/GoDump/internal/service/drive"
	"github.com/BrunoTulio/GoDump/internal/service/process_dump"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type (
	RestoreAutoUseCase interface {
		Restore(enableDriver bool) error
	}

	restoreAutoUseCase struct {
		configurationFileService configuration_file.ConfigurationFileService
		backupFilesService       backup_files.BackupFilesService
		processDumpService       process_dump.ProcessDumpService
		driveService             drive.DriveService
	}
)

// Restore implements RestoreAutoUseCase.
func (r *restoreAutoUseCase) Restore(enableDriver bool) error {
	config, err := r.configurationFileService.Get()

	if err != nil {
		logger.Fatal(err)
	}

	err = config.IsValid()

	if err != nil {
		logger.Fatal(err)
	}

	if enableDriver {
		logger.Info("Enable restore driver")
		return r.restoreAllDriver(config)
	}

	return r.restoreAllLocalFile(config)

}

func (r *restoreAutoUseCase) restoreAllLocalFile(config domain.Config) error {
	sort.Sort(config.Dumps)
	group := errgroup.Group{}
	for _, dump := range config.Dumps {
		dump := dump
		group.Go(func() error {
			err := r.restoreLocalFile(dump)
			if err != nil {
				logger.Error(err)
				return err
			}
			return nil
		})
	}
	return group.Wait()
}

func (r *restoreAutoUseCase) restoreAllDriver(config domain.Config) error {
	sort.Sort(config.Dumps)
	group := errgroup.Group{}
	for _, dump := range config.Dumps {
		dump := dump
		group.Go(func() error {
			err := r.restoreDrive(dump)
			if err != nil {
				logger.Error(err)
				return err
			}
			return nil
		})
	}
	return group.Wait()
}

func (r *restoreAutoUseCase) restoreLocalFile(dump domain.Dump) error {

	file, err := r.backupFilesService.GetLast(dump)

	if err != nil {
		return err
	}
	logger.Infof("Last file backup: %v", file.Path)
	err = r.processDumpService.Restore(file.Path, dump)

	if err != nil {
		return err
	}

	return nil
}

func (r *restoreAutoUseCase) restoreDrive(dump domain.Dump) error {

	path, err := r.driveService.Download(context.Background(), dump)

	if err != nil {
		return err
	}

	err = r.processDumpService.Restore(path, dump)

	os.Remove(path)

	if err != nil {
		return err
	}

	return nil
}

func NewRestoreAuto(
	configurationFileService configuration_file.ConfigurationFileService,
	backupFilesService backup_files.BackupFilesService,
	processDumpService process_dump.ProcessDumpService,
	driveService drive.DriveService,
) RestoreAutoUseCase {
	return &restoreAutoUseCase{
		configurationFileService,
		backupFilesService,
		processDumpService,
		driveService,
	}
}
