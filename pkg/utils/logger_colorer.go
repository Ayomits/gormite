package utils

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func NewLogger(module string) *log.Logger {
	opts := log.Options{
		ReportTimestamp: true,
		Prefix:          module,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
	}
	logger := log.NewWithOptions(os.Stdout, opts)
	return logger
}
