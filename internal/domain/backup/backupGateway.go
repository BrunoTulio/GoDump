package backup

import (
	"context"
	"time"
)

type (
	BackupRepository interface {
		StoreLastMail(time.Time) error
		StoreLastFile(time.Time) error
		LastFile() (*time.Time, error)
		LastMail() (*time.Time, error)
		First() (*Backup, error)
	}

	BackupMake interface {
		Execute(ctx context.Context, nameFile string) (string, error)
	}

	BackupCheckConnection interface {
		Execute(ctx context.Context) error
	}

	BackupCheckDump interface {
		Execute() error
	}
)
