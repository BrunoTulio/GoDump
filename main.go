package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/BrunoTulio/GoDump/internal/applications/backup/listener"
	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/internal/infra/backup/adapter"
	"github.com/BrunoTulio/GoDump/internal/infra/backup/repository"
	"github.com/BrunoTulio/GoDump/internal/infra/http"
	"github.com/BrunoTulio/GoDump/internal/infra/mail/simplemail"
	badger "github.com/dgraph-io/badger/v3"

	"github.com/BrunoTulio/GoDump/internal/settings"
	"github.com/BrunoTulio/GoDump/pkg/event"
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/BrunoTulio/GoDump/pkg/logger/zap"

	"github.com/robfig/cron/v3"
)

func main() {
	setup()
}

func setup() {

	settings.Load()

	log := logger.NewLogger(zap.NewZapLogger())

	backupCheckDumb := adapter.NewBackupCheckDumb()

	backupMake := adapter.NewBackupMake()

	backupCheckConnection := adapter.NewBackupCheckConnection()

	mailSmtp := simplemail.New(
		settings.MailHost(),
		settings.MailUser(),
		settings.MailPassword(),
		settings.MailPort(),
		settings.MailSecure(),
	)

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	backupRepository := repository.New(db)

	backupListener := listener.NewBackupListener(
		backupCheckDumb,
		backupCheckConnection,
		backupMake,
		backupRepository,
		mailSmtp,
	)

	event.On(backup.BackupEventKey,
		backupListener,
		event.Normal,
	)

	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger), // or use cron.DefaultLogger
	))
	defer c.Stop()

	c.Schedule(cron.Every(settings.JobIntervalFile()), cron.FuncJob(func() {
		log.Info("Initialize event backup database file")
		err := event.AwaitFire(context.Background(), backup.NewBackupLocalEvent())

		if err != nil {
			log.Errorf("Backup failed: %s", err)
		}
	}))

	c.Schedule(cron.Every(settings.JobIntervalMail()), cron.FuncJob(func() {
		log.Info("Initialize event backup database send mail")
		err := event.AwaitFire(context.Background(), backup.NewBackupMailEvent())

		if err != nil {
			log.Errorf("Backup failed send mail: %s", err)
		}
	}))

	log.Info("Initializer schedule backup postgres")

	go c.Start()
	go http.Setup(backupRepository)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
	log.Info(fmt.Sprintf("%s closed", time.Now().Format("2006-01-02 15:04:05")))
}
