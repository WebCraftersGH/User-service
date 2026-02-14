package logging

import (
	"log/slog"
	"os"
)

type logger struct {
	l *slog.Logger
}

func NewLogger(logLevel string) *logger {

	var lLevel slog.Leveler

	switch logLevel {
	case "DEBUG":
		lLevel = slog.LevelDebug
	case "WARN":
		lLevel = slog.LevelWarn
	case "ERROR":
		lLevel = slog.LevelError
	default:
		lLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lLevel,
	})
	return &logger{
		l: slog.New(handler),
	}
}

func (lg *logger) Info(msg string, args ...any) {
	lg.l.Info(msg, args...)
}

func (lg *logger) Debug(msg string, args ...any) {
	lg.l.Debug(msg, args...)
}

func (lg *logger) Error(msg string, args ...any) {
	lg.l.Error(msg, args...)
}

func (lg *logger) Warn(msg string, args ...any) {
	lg.l.Warn(msg, args...)
}
