package factory

import (
	"github.com/BrunoTulio/GoDump/pkg/logger"
	"github.com/BrunoTulio/GoDump/pkg/logger/zap.v1"
)

func MakeLogger() logger.Logger {
	return zap.NewWithOptions(zap.NewOptions(
		struct {
			Enabled   bool
			Level     string
			Formatter string
		}{
			Enabled:   true,
			Level:     "DEBUG",
			Formatter: "TEXT",
		},
		struct {
			Enabled   bool
			Level     string
			Path      string
			Name      string
			MaxSize   int
			Compress  bool
			MaxAge    int
			Formatter string
		}{
			Enabled: false,
			Level:   "DEBUG",
			Path:    "./",
			Name:    "log",
		},
	))

}
