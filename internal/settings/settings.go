package settings

import "time"

type Setting struct {
	databaseHost     string
	databasePort     int
	databaseUser     string
	databasePassword string
	databaseName     string
	databaseSSL      bool

	pathFile string

	intervalFile time.Duration
	intervalMail time.Duration

	mailHost     string
	mailPort     int
	mailUser     string
	mailPassword string
	mailSecure   bool
	mailSend     string
}

func (s *Setting) DatabaseHost() string {
	return s.databaseHost
}

func (s *Setting) DatabasePort() int {
	return s.databasePort
}

func (s *Setting) DatabaseUser() string {
	return s.databaseUser
}

func (s *Setting) DatabasePassword() string {
	return s.databasePassword
}

func (s *Setting) DatabaseName() string {
	return s.databaseName
}

func (s *Setting) DatabaseSSL() bool {
	return s.databaseSSL
}

func (s *Setting) IntervalFile() time.Duration {
	return s.intervalFile
}

func (s *Setting) IntervalSendMail() time.Duration {
	return s.intervalMail
}

func (s *Setting) MailHost() string {
	return s.mailHost
}

func (s *Setting) MailPort() int {
	return s.mailPort
}

func (s *Setting) MailUser() string {
	return s.mailUser
}

func (s *Setting) MailPassword() string {
	return s.mailPassword
}

func (s *Setting) MailSecure() bool {
	return s.mailSecure
}

func (s *Setting) MailSend() string {
	return s.mailSend
}

func (s *Setting) PathFile() string {
	return s.pathFile
}
