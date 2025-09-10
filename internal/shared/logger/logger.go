package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func NewLogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	},
	)

	return slog.New(handler)
}

// NewFileLogger writes JSON logs to stdout + a rotating file
// named app-YYYY-MM-DD.log under dir.
func NewFileLogger(level slog.Level, dir string, maxSizeMB int) *slog.Logger {
	if dir == "" {
		// Fallback to stdout-only if misconfigured
		return NewLogger(level)
	}

	filename := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	absPath := filepath.Join(dir, filename)

	// Ensure directory exists
	_ = os.MkdirAll(filepath.Dir(absPath), 0o755)

	rotator := &lumberjack.Logger{
		Filename:   absPath,
		MaxSize:    maxSizeMB,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	out := io.MultiWriter(os.Stdout, rotator)

	h := slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})
	l := slog.New(h)

	l.Info("logger initialized", "file", absPath)

	return l
}
