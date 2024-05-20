package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func ellog() {
	w := os.Stderr

	// create a new logger
	logger := slog.New(tint.NewHandler(w, nil))
	slog.SetDefault(logger)

	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}
