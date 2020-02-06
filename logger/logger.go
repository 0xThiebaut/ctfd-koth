package logger

import (
	"log"
	"os"
)

// Let's be lazy and have some predefined loggers.
var (
	// Info logs any informational event to STDOUT.
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags)

	// Warn flags non-critical events in STDERR.
	Warn = log.New(os.Stderr, "[WARN] ", log.LstdFlags)

	// Critical flags critical events in STDERR.
	Critical = log.New(os.Stderr, "[CRITICAL] ", log.LstdFlags)
)
