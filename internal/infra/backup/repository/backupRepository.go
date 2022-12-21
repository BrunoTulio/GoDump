package repository

import (
	"time"

	"github.com/BrunoTulio/GoDump/internal/domain/backup"

	badger "github.com/dgraph-io/badger/v3"
)

const keyLastFile = "backup.last-file"
const keyLastMail = "backup.last-mail"

type backupRepository struct {
	db *badger.DB
}

// First implements backup.BackupRepository
func (b backupRepository) First() (*backup.Backup, error) {
	lastFile, err := b.LastMail()

	if err != nil {
		return nil, err
	}

	lastMail, err := b.LastFile()

	if err != nil {
		return nil, err
	}

	return &backup.Backup{
		LastDateFile:     lastFile,
		LastDateSendMail: lastMail,
	}, nil
}

// LastFile implements backup.BackupRepository
func (b backupRepository) LastFile() (*time.Time, error) {
	var last *time.Time
	err := b.db.View(func(txn *badger.Txn) error {
		iten, err := txn.Get([]byte(keyLastFile))

		if err != nil {
			return err
		}

		err = iten.Value(func(val []byte) error {
			b, err := time.Parse(time.RFC3339, string(val))

			if err != nil {
				return err
			}
			last = &b
			return nil
		})

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return last, nil

}

// LastMail implements backup.BackupRepository
func (b backupRepository) LastMail() (*time.Time, error) {
	var last *time.Time
	err := b.db.View(func(txn *badger.Txn) error {
		iten, err := txn.Get([]byte(keyLastMail))

		if err != nil {
			return err
		}

		err = iten.Value(func(val []byte) error {
			b, err := time.Parse(time.RFC3339, string(val))

			if err != nil {
				return err
			}
			last = &b
			return nil
		})

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return last, nil
}

// StoreLastFile implements backup.BackupRepository
func (b backupRepository) StoreLastFile(t time.Time) error {
	return b.db.Update(func(txn *badger.Txn) error {
		value := t.Format(time.RFC3339)
		return txn.Set([]byte(keyLastFile), []byte(value))
	})
}

// StoreLastMail implements backup.BackupRepository
func (b backupRepository) StoreLastMail(t time.Time) error {
	return b.db.Update(func(txn *badger.Txn) error {
		value := t.Format(time.RFC3339)
		return txn.Set([]byte(keyLastMail), []byte(value))
	})
}

func New(db *badger.DB) backup.BackupRepository {

	return backupRepository{
		db: db,
	}
}
