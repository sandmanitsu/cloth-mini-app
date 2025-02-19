package logger

import (
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

// Create slog logger
func NewLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("logger error: env is empty or incorrect value")
	}

	return logger
}

// Wrap error to slog.Attr struct
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
