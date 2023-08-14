package zap

type (
	Options struct {
		Console struct {
			Enabled   bool   // enable/disable console logging
			Level     string // console log level
			Formatter string // console formatter TEXT/JSON
		}
		File struct {
			Enabled   bool   // enable/disable file logging
			Level     string // file log level
			Path      string // file log path
			Name      string // file log filename
			MaxSize   int    // log file max size (MB)
			Compress  bool   // enabled/disable file compress
			MaxAge    int    // file max age
			Formatter string // file formatter TEXT/JSON
		}
	}
)

func NewOptions(console struct {
	Enabled   bool
	Level     string
	Formatter string
}, file struct {
	Enabled   bool
	Level     string
	Path      string
	Name      string
	MaxSize   int
	Compress  bool
	MaxAge    int
	Formatter string
}) *Options {
	return &Options{
		Console: console,
		File:    file,
	}
}
