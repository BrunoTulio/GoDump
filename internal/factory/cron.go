package factory

import "github.com/robfig/cron/v3"

func MakeCron() *cron.Cron {
	return cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger), // or use cron.DefaultLogger
	))
}
