package adapter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/internal/settings"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/pkg/errors"
)

var _ backup.BackupMake = backupMakeAdapter{}

type backupMakeAdapter struct {
}

// Execute implements backup.BackupDumbProcessAdapter
func (b backupMakeAdapter) Execute(ctx context.Context, nameFile string) (string, error) {

	pathFile := path.Join(settings.PathFile(), fmt.Sprintf("%s.sql.tar.gz", nameFile))

	cmd := exec.Command(
		"pg_dump",
		"--host",
		settings.DatabaseHost(),
		"--port",
		strconv.FormatInt(int64(settings.DatabasePort()), 10),
		"--username",
		settings.DatabaseUser(),
		"--dbname",
		settings.DatabaseName(),
		"--format",
		"c",
		"--file",
		pathFile,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSWORD=%s`, settings.DatabasePassword()))
	stderrIn, err := cmd.StderrPipe()

	if err != nil {
		return "", errors.Wrap(err, "failed to dump")
	}

	defer stderrIn.Close()

	output := ""
	go func() {
		reader := bufio.NewReader(stderrIn)
		line, err := reader.ReadString('\n')

		output += line
		for err == nil {
			line, err = reader.ReadString('\n')
			output += line
		}
	}()

	err = cmd.Start()

	if err != nil {
		return "", errors.Wrap(err, "failed to dump")
	}

	err = cmd.Wait()

	if err != nil {
		return "", errors.Wrap(err, "failed to dump")
	}

	logger.Infof("Dump finish: %s Output: %s\n", pathFile, output)

	return pathFile, nil

}

func NewBackupMake() backup.BackupMake {
	return backupMakeAdapter{}
}
