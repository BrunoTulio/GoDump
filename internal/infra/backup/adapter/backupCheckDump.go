package adapter

import (
	"os/exec"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
)

var _ backup.BackupCheckDump = backupCheckDump{}

type backupCheckDump struct {
}

func NewBackupCheckDumb() backup.BackupCheckDump {
	return backupCheckDump{}
}

// Execute implements backup.BackupCheckDumpAdapter
func (b backupCheckDump) Execute() error {
	_, err := exec.LookPath("pg_dump")
	if err != nil {
		return err
	}
	return nil
}
