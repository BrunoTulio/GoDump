package backup

import "github.com/BrunoTulio/GoDump/pkg/event"

const BackupEventKey = "backup.event"

type BackupLocalEvent struct {
}

// Process implements ProcessEventStrategy
func (b *BackupLocalEvent) Event() event.Event {
	return b
}

func (b *BackupLocalEvent) Key() string {
	return BackupEventKey
}

func NewBackupLocalEvent() event.Event {
	return &BackupLocalEvent{}
}

type BackupMailEvent struct {
}

// Process implements ProcessEventStrategy
func (b *BackupMailEvent) Event() event.Event {
	return b
}

func (b *BackupMailEvent) Key() string {
	return BackupEventKey
}

func NewBackupMailEvent() event.Event {
	return &BackupMailEvent{}
}
