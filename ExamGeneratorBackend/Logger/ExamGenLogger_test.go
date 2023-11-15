package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLogger(t *testing.T) {
	config := LoggerConfig{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJson:      true,
		FileLoggingEnabled:    true,
		Directory:             "log",
		Filename:              "test.log",
		MaxSize:               10,
		MaxBackups:            3,
		MaxAge:                7,
	}

	logger := ConfigLogger(config)
	if logger == nil {
		t.Fatalf("Failed to create logger")
	}
}

func TestNewRollingFile(t *testing.T) {
	config := LoggerConfig{
		Directory:  "log",
		Filename:   "test.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	}

	writer := newRollingFile(config)
	if writer == nil {
		t.Fatalf("Failed to create writer")
	}

	// Clean up
	os.Remove(filepath.Join(config.Directory, config.Filename))
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if logger == nil {
		t.Fatalf("Failed to get logger")
	}
}
