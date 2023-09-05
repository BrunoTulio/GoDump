package backup

import (
	"fmt"
	"sort"

	"github.com/BrunoTulio/GoDump/internal/domain"

	"github.com/BrunoTulio/GoDump/internal/service/configuration_file"
	"github.com/BrunoTulio/GoDump/internal/service/process_dump"
	"golang.org/x/sync/errgroup"
)

type (
	BackupUseCase interface {
		Generate(fileType *domain.Type) error
		GenerateByKey(fileType *domain.Type, key string) error
	}
	backupUseCase struct {
		configurationFileService configuration_file.ConfigurationFileService
		processDumpService       process_dump.ProcessDumpService
	}
)

// Generate implements BackupUseCase.
func (b *backupUseCase) Generate(fileType *domain.Type) error {

	config, err := b.configurationFileService.Get()

	if err != nil {
		return err
	}

	err = config.IsValid()

	if err != nil {
		return err
	}

	if fileType != nil {
		config.Dumps.UpdateType(fileType)
	}

	sort.Sort(config.Dumps)
	group := errgroup.Group{}

	for _, dump := range config.Dumps {
		dump := dump
		group.Go(func() error {
			_, err := b.processDumpService.Backup(dump)
			if err != nil {
				return err
			}
			return nil
		})

	}

	return group.Wait()

}

// GenerateByKey implements BackupUseCase.
func (b *backupUseCase) GenerateByKey(fileType *domain.Type, key string) error {

	config, err := b.configurationFileService.Get()

	if err != nil {
		return err
	}

	err = config.IsValid()

	if err != nil {
		return err
	}

	var dump domain.Dump

	for _, d := range config.Dumps {
		if key == d.Key {
			dump = d
		}
	}

	if dump.Key == "" {
		return fmt.Errorf("Dump key [%s] not found", key)
	}
	group := errgroup.Group{}
	group.Go(func() error {
		_, err := b.processDumpService.Backup(dump)
		if err != nil {
			return err
		}
		return nil
	})

	return group.Wait()

}

func NewBackup(
	configurationFileService configuration_file.ConfigurationFileService,
	processDumpService process_dump.ProcessDumpService,
) BackupUseCase {
	return &backupUseCase{
		configurationFileService,
		processDumpService,
	}
}
