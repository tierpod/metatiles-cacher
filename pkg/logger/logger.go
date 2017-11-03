// Package logger helps to configure logger with filtering.
package logger

import (
	"io"
	"log"

	"github.com/hashicorp/logutils"
)

// New returns new logger with filtered levels:
//   {"DEBUG", "WARN", "ERROR", "INFO"}
//
// If debug=true, show messages from "DEBUG" level.
//
// If datetime=true, show datetime in log message.
func New(out io.Writer, debug, datetime bool) *log.Logger {
	levels := []logutils.LogLevel{"DEBUG", "WARN", "ERROR", "INFO"}

	var level string
	switch debug {
	case true:
		level = "DEBUG"
	default:
		level = "WARN"
	}

	var flag int
	switch datetime {
	case false:
		flag = 0
	default:
		flag = log.LstdFlags
	}

	logger := log.New(out, "", flag)
	filter := &logutils.LevelFilter{
		Levels:   levels,
		MinLevel: logutils.LogLevel(level),
		Writer:   out,
	}

	logger.SetOutput(filter)

	return logger
}
