package process_dump

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/BrunoTulio/GoDump/internal/constants"
	"github.com/BrunoTulio/GoDump/internal/domain"
	"github.com/BrunoTulio/GoDump/pkg/folder"
	"github.com/BrunoTulio/GoDump/pkg/logger"

	"github.com/pkg/errors"
)

type (
	ProcessDumpService interface {
		Backup(dump domain.Dump) (string, error)
		Restore(pathFile string, dump domain.Dump) error
	}

	processDumpService struct {
	}
)

// Restore implements ProcessDumpService.
func (s *processDumpService) Restore(pathFile string, dump domain.Dump) (err error) {
	cmd := exec.Command(
		"pg_restore",
		"--host",
		dump.Host, //container db host
		"--port",
		fmt.Sprintf("%d", dump.Port), //container db port
		"--username",
		dump.Username, //container db username
		"--dbname",
		dump.Database, //container db database_name

		pathFile, //container loca/arquivo backup
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSWORD=%s`, dump.Password))

	stderrIn, err := cmd.StderrPipe()

	if err != nil {
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
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
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
	}

	err = cmd.Wait()

	if err != nil {
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
	}

	logger.Infof("Restore %s finish: %s Output: %s\n", dump.Key, pathFile, output)

	return

}

// Backup implements StartCronBackup.
func (s *processDumpService) Backup(dump domain.Dump) (path string, err error) {

	path = pathFile(dump)

	args := []string{
		"--host",
		dump.Host, //container db host
		"--port",
		fmt.Sprintf("%d", dump.Port), //container db port
		"--username",
		dump.Username, //container db username
		"--dbname",
		dump.Database,
		"--file",
		path, //container loca/arquivo backup
	}

	if !dump.IsTypeSQL() {
		args = append(args,
			"--format",
			"tar",
		)
	}

	cmd := exec.Command(
		"pg_dump",
		args...,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSWORD=%s`, dump.Password))

	stderrIn, err := cmd.StderrPipe()

	if err != nil {
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
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
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
	}

	err = cmd.Wait()

	if err != nil {
		err = errors.Wrapf(err, "Dump %s failed", dump.Key)
		return
	}

	logger.Infof("Dump %s finish: %s Output: %s\n", dump.Key, path, output)

	return
}

func fileName(dump domain.Dump) string {
	date := time.Now()
	fileName := fmt.Sprintf("%s_%s%s",
		dump.Database,
		date.Format(constants.LayoutDate),
		dump.Extension(),
	)

	return fileName
}

func pathFile(dump domain.Dump) string {
	err := folder.Create(constants.PathBackup)

	if err != nil {
		logger.Fatal(err)
	}

	dir := path.Join(constants.PathBackup, dump.Key)

	err = folder.Create(dir)

	if err != nil {
		logger.Fatal(err)
	}

	return path.Join(dir, fileName(dump))
}

func New() ProcessDumpService {
	return &processDumpService{}
}
