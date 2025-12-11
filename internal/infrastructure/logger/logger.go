package logger

import (
	"io"
	"log"
	"os"
)

// Logger is a thin wrapper around log.Logger to support simple
// structured-style logging. For this exercise we keep it minimal,
// but the abstraction lets us swap implementation later (Elasticsearch, zerolog, etc.).
type Logger struct {
	l *log.Logger
}

// New creates a Logger that writes to stdout with standard flags.
func New() *Logger {
	return NewWithWriter(os.Stdout)
}

// NewWithWriter creates a Logger that writes to the provided writer.
// This is mainly useful for tests.
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{
		l: log.New(w, "", log.LstdFlags),
	}
}

// Info logs an informational message with key/value pairs.
func (lg *Logger) Info(msg string, kv ...any) {
	args := append([]any{"level=INFO", msg}, kv...)
	lg.l.Println(args...)
}

// Error logs an error message with key/value pairs.
func (lg *Logger) Error(msg string, kv ...any) {
	args := append([]any{"level=ERROR", msg}, kv...)
	lg.l.Println(args...)
}
