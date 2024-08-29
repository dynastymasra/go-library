package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZeroLogConfig is a struct that holds the configuration settings for Zero log.
// Zero log is a zero allocation JSON logger in Go that is fast and reliable.
//
// The struct fields are:
// Level: The logging level. It can be "debug", "info", "warn", "error", "fatal", or "panic".
// FileEnabled: A boolean indicating whether logging to a file is enabled.
// FilePath: The path to the log file. This is used if FileEnabled is true.
// FileMaxSize: The maximum size of the log file in megabytes.
// FileMaxBackup: The maximum number of old log files to retain.
// FileMaxAge: The maximum number of days to retain old log files.
type ZeroLogConfig struct {
	Level         string
	FileEnabled   bool // FileEnabled whether logging to a file is enabled.
	FilePath      string
	FileMaxSize   int // FileMaxSize the maximum size of the log file.
	FileMaxBackup int // FileMaxBackup the maximum number of old log files to retain.
	FileMaxAge    int // FileMaxAge the maximum number of days to retain old log files.
}

// ConfigureZeroLog sets up the ZeroLog logger based on the provided configuration.
// It configures the global logging level, sets up the log output destinations (console and/or file),
// and adds contextual information such as hostname, service name, and version to each log entry.
//
// Parameters:
// - name: The name of the service or application.
// - version: The version of the service or application.
func (z ZeroLogConfig) ConfigureZeroLog(name, version string) {
	switch strings.ToLower(z.Level) {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var writers []io.Writer
	if z.FileEnabled {
		// Output is JSON.
		writers = append(writers, &lumberjack.Logger{
			Filename:   z.FilePath,
			MaxSize:    z.FileMaxSize,
			MaxAge:     z.FileMaxAge,
			MaxBackups: z.FileMaxBackup,
		})
	} else {
		writers = append(writers, zerolog.NewConsoleWriter())
	}

	mw := io.MultiWriter(writers...)
	log.Logger = zerolog.New(mw).With().Timestamp().Caller().
		Str("hostname", hostname).Str("service", name).Str("version", version).Logger()
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = time.Now().UTC
}

// InitializeZeroLogTestLogger is a function that sets up a logger for testing purposes.
// It disables the global logging level and discards all log output.
// This is useful in a testing context where log output might not be desired.
func InitializeZeroLogTestLogger() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	log.Logger = zerolog.New(io.Discard)
}
