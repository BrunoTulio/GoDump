package backup

import "time"

type Backup struct {
	LastDateFile     *time.Time
	LastDateSendMail *time.Time
}

func TimeFormat(t *time.Time) *string {
	if t == nil {
		return nil
	}

	v := t.Format(time.RFC3339)
	return &v
}

func (b *Backup) LastDateFileFormat() *string {
	return TimeFormat(b.LastDateFile)
}

func (b *Backup) LastDateMailFormat() *string {
	return TimeFormat(b.LastDateSendMail)
}
