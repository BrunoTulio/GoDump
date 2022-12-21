package settings

import (
	"os"
	"time"

	"github.com/BrunoTulio/GoDump/pkg/env"
)

const (
	root = "godump"

	databaseHost     = root + ".db.host"
	databasePort     = root + ".db.port"
	databaseUser     = root + ".db.user"
	databasePassword = root + ".db.password"
	databaseName     = root + ".db.name"
	databaseSSL      = root + ".db.ssl"

	mailHost     = root + ".mail.host"
	mailPort     = root + ".mail.port"
	mailUser     = root + ".mail.user"
	mailPassword = root + ".mail.password"
	mailSecure   = root + ".mail.secure"
	mailSend     = root + ".mail.send"

	pathFile = root + ".path-file"

	jobIntervalFile = root + ".job.interval.file"
	jobIntervalMail = root + ".job.interval.mail"
)

var s *Setting

func Load() {

	// currentPath, _ := os.Getwd()

	// env.EnvLoadFiles(path.Join(currentPath, ".env"))

	pathFile := envPathFile()
	databaseHost := envDatabaseHost()
	databasePort := envDatabasePort()
	databaseUser := envDatabaseUser()
	databasePassword := envDatabasePassword()
	databaseName := envDatabaseName()
	databaseSSL := env.BoolDefault(databaseSSL, false)

	s = &Setting{
		databaseHost:     databaseHost,
		databasePort:     databasePort,
		databaseUser:     databaseUser,
		databasePassword: databasePassword,
		databaseName:     databaseName,
		databaseSSL:      databaseSSL,

		mailHost:     env.StringDefault(mailHost, ""),
		mailUser:     env.StringDefault(mailUser, ""),
		mailPassword: env.StringDefault(mailPassword, ""),
		mailSecure:   env.BoolDefault(mailSecure, false),
		mailPort:     env.IntDefault(mailPort, 385),
		mailSend:     envMailSend(),

		pathFile: pathFile,

		intervalFile: env.DurationDefault(jobIntervalFile, time.Minute*1),
		intervalMail: env.DurationDefault(jobIntervalMail, time.Minute*2),
	}

}

func envMailSend() string {
	mailSend, err := env.MustString(mailSend)

	if err != nil {
		panic(err)
	}
	return mailSend
}

func envDatabasePort() int {
	databasePort, err := env.MustInt(databasePort)
	if err != nil {
		panic(err)
	}
	return databasePort
}

func envDatabaseHost() string {
	databaseHost, err := env.MustString(databaseHost)

	if err != nil {
		panic(err)
	}
	return databaseHost
}

func envDatabaseUser() string {
	databaseUser, err := env.MustString(databaseUser)

	if err != nil {
		panic(err)
	}
	return databaseUser
}

func envDatabaseName() string {
	databaseName, err := env.MustString(databaseName)

	if err != nil {
		panic(err)
	}
	return databaseName
}

func envDatabasePassword() string {
	databasePassword, err := env.MustString(databasePassword)

	if err != nil {
		panic(err)
	}
	return databasePassword
}

func envPathFile() string {
	pathFile := env.String(pathFile)

	if pathFile == "" {
		path, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		pathFile = path
	}
	return pathFile
}

func DatabaseHost() string {
	return s.databaseHost
}

func DatabasePort() int {
	return s.databasePort
}

func DatabaseUser() string {
	return s.databaseUser
}

func DatabasePassword() string {
	return s.databasePassword
}

func DatabaseName() string {
	return s.databaseName
}

func DatabaseSSL() bool {
	return s.databaseSSL
}

func MailHost() string {
	return s.mailHost
}

func MailPort() int {
	return s.mailPort
}

func MailUser() string {
	return s.mailUser
}

func MailPassword() string {
	return s.mailPassword
}

func MailSecure() bool {
	return s.mailSecure
}

func MailSend() string {
	return s.mailSend
}

func JobIntervalFile() time.Duration {
	return s.intervalFile
}

func JobIntervalMail() time.Duration {
	return s.intervalMail
}

func PathFile() string {
	return s.pathFile
}

func Settings() Setting {
	return *s
}

func DatabaseSSLDescription() string {
	if s.databaseSSL {
		return "enable"
	}
	return "disable"
}
