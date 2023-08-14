package backup_files

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/pkg/errors"
)

type (
	BackupFile struct {
		Date time.Time
		Diff time.Duration
		Name string
		Path string
	}

	BackupFiles []BackupFile

	BackupFilesService interface {
		RemoveByDuration(diff time.Duration, dump domain.Dump) error
		GetLast(dump domain.Dump) (BackupFile, error)
	}

	backupFilesService struct{}
)

// GetLast implements BackupFilesService.
func (r *backupFilesService) GetLast(dump domain.Dump) (BackupFile, error) {
	files, err := readFiles(dump)

	if err != nil {
		return BackupFile{}, errors.Wrap(err, "ReadFiles failed")
	}

	sort.Sort(sort.Reverse(files))

	if len(files) == 0 {
		return BackupFile{}, fmt.Errorf("Not found backup las key %s", dump.Key)
	}

	return files[0], nil
}

// RemoveByDuration implements RemoveBackupsService.
func (r *backupFilesService) RemoveByDuration(diff time.Duration, dump domain.Dump) error {

	files, err := readFiles(dump)

	if err != nil {
		return errors.Wrap(err, "ReadFiles failed")
	}

	var filesRemove []string

	for _, f := range files {
		if f.Diff >= diff {
			filesRemove = append(filesRemove, f.Path)
		}
	}

	for _, f := range filesRemove {
		err = os.RemoveAll(f)

		if err != nil {
			return errors.Wrapf(err, "Error removing file %s", f)
		}

		logger.Infof("Remove %s\n", f)
	}

	return nil

}

func readFiles(dump domain.Dump) (BackupFiles, error) {

	var backupFiles BackupFiles

	dir := path.Join(constants.PathBackup, dump.Key)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDir failed")
	}

	for _, d := range entries {
		file, err := d.Info()

		if err != nil {
			return nil, errors.Wrap(err, "Error reading file")
		}

		name := file.Name()

		date, err := filaNameToDate(name)

		if err != nil {
			return nil, err
		}

		diff := time.Now().Sub(date)

		backupFiles = append(backupFiles, BackupFile{
			Date: date,
			Diff: diff,
			Name: name,
			Path: path.Join(dir, name),
		})
	}

	return backupFiles, nil
}

func filaNameToDate(nameFile string) (time.Time, error) {
	name := strings.ReplaceAll(nameFile, ".dump", "")
	values := strings.Split(name, "_")

	if len(values) != 2 {
		return time.Time{}, fmt.Errorf("Invalid name file")
	}

	date, err := time.Parse(constants.LayoutDate, values[1])

	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid date")
	}

	return date, nil
}

func New() BackupFilesService {
	return &backupFilesService{}
}

func (c BackupFiles) Len() int {
	return len(c)
}
func (c BackupFiles) Less(i, j int) bool {
	return c[i].Date.Before(c[j].Date)
}
func (c BackupFiles) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
