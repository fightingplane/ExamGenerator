package logging

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration for logging
type LoggerConfig struct {

	// Enable console logging
	ConsoleLoggingEnabled bool `json:"consoleLoggingEnabled" yaml:"consoleLoggingEnabled"`

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool `json:"encodeLogsAsJson" yaml:"encodeLogsAsJson"`

	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `json:"filenameEnabled" yaml:"filenameEnabled"`

	// Directory to log to to when filelogging is enabled
	Directory string `json:"directory" yaml:"directory"`

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// MaxAge the max age in days to keep a logfile
	MaxAge int `json:"maxage" yaml:"maxage"`
}

type ExamGenLogger struct {
	*zerolog.Logger
}

var (
	LoggerInstance *ExamGenLogger
	lock           = &sync.Mutex{}
)

func ConfigLogger(config LoggerConfig) *ExamGenLogger {

	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	return &ExamGenLogger{
		Logger: &logger,
	}
}

func newRollingFile(config LoggerConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}

func GetLogger() *ExamGenLogger {

	if LoggerInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if LoggerInstance == nil {
			config := &LoggerConfig{}
			viper.UnmarshalKey("LoggerConfigurations", config)
			LoggerInstance = ConfigLogger(*config)
		} else {
			LoggerInstance.Info().Msg("Logger already initialized")
		}
	}

	return LoggerInstance
}
