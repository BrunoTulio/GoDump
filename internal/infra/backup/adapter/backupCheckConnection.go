package adapter

import (
	"context"
	"fmt"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/internal/settings"

	"github.com/jackc/pgx/v5"
)

type backupCheckConnection struct {
}

// Execute implements backup.BackupCheckConnection
func (b backupCheckConnection) Execute(ctx context.Context) (err error) {

	//postgres://proxy:proxy@db:5432/proxy?sslmode=disable
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		settings.DatabaseUser(),
		settings.DatabasePassword(),
		settings.DatabaseHost(),
		settings.DatabasePort(),
		settings.DatabaseName(),
		settings.DatabaseSSLDescription(),
	)

	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return
	}
	defer func() { err = conn.Close(context.Background()) }()

	err = conn.Ping(ctx)

	if err != nil {
		return
	}

	return
}

func NewBackupCheckConnection() backup.BackupCheckConnection {
	return backupCheckConnection{}
}
