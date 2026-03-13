package logger

import (
	"log/slog"
	"os"

	"github.com/identicalaffiliation/orders-procceser-api/migrator/internal/config"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type slogger struct {
	logger *slog.Logger
}

func NewLogger(cfg *config.MigratorConfig) Logger {
	var l slog.Level
	switch cfg.LoggerConfig.Level {
	case "debug":
		l = slog.LevelDebug
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: l,
	}

	h := slog.NewJSONHandler(os.Stdout, opts)

	return &slogger{
		logger: slog.New(h),
	}
}

func (l *slogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *slogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
