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

// ConfigureZeroLog configures the global logger with the settings from the ZeroLogConfig.
// It sets the global logging level, hostname, service name, and version.
// If FileEnabled is true, it also sets up a file logger.
//
// Parameters:
// name: The name of the service. This will be included in every log entry.
// version: The version of the service. This will be included in every log entry.
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
	writers = append(writers, zerolog.NewConsoleWriter())

	if z.FileEnabled {
		// Output is JSON.
		writers = append(writers, &lumberjack.Logger{
			Filename:   z.FilePath,
			MaxSize:    z.FileMaxSize,
			MaxAge:     z.FileMaxAge,
			MaxBackups: z.FileMaxBackup,
		})
	}

	mw := io.MultiWriter(writers...)
	log.Logger = zerolog.New(mw).With().Timestamp().Caller().
		Str("hostname", hostname).Str("service", name).Str("version", version).Logger()
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

// InitializeZeroLogTestLogger is a function that sets up a logger for testing purposes.
// It disables the global logging level and discards all log output.
// This is useful in a testing context where log output might not be desired.
func InitializeZeroLogTestLogger() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	log.Logger = zerolog.New(io.Discard)
}
