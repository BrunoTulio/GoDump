package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"
	"github.com/BrunoTulio/GoDump/pkg/logger"
)

func Last(repository backup.BackupRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		backup, err := repository.First()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})

			if err != nil {
				logger.Error(err)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"lastFile": backup.LastDateFileFormat(),
			"lastMail": backup.LastDateMailFormat(),
		})

		if err != nil {
			logger.Error(err)
		}

	}
}
