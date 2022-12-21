package http

import (
	"net/http"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/internal/infra/http/handler"
	"github.com/BrunoTulio/GoDump/pkg/logger"
)

func Setup(backupRepository backup.BackupRepository) {

	http.HandleFunc("/last/file", handler.LastFile(backupRepository))
	http.HandleFunc("/last/mail", handler.LastMail(backupRepository))
	http.HandleFunc("/last", handler.Last(backupRepository))

	logger.Fatal(http.ListenAndServe(":5600", nil))
}
