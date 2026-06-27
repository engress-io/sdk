package logx

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

var savedDefault *slog.Logger

// Init configures the default slog logger from FLUX_LOG_FORMAT (json|text).
// Service logs go to stderr so the agent TUI can own stdout.
func Init() {
	InitWriter(os.Stderr)
}

// InitWriter configures slog using w (tests).
func InitWriter(w io.Writer) {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("FLUX_LOG_FORMAT")))
	opts := log.Options{
		ReportTimestamp: true,
		Level:           log.InfoLevel,
	}
	if format == "json" {
		opts.Formatter = log.JSONFormatter
	}
	logger := log.NewWithOptions(w, opts)
	slog.SetDefault(slog.New(logger))
}

// Quiet discards slog output. Use when a TUI owns the terminal (stderr would corrupt the display).
func Quiet() {
	savedDefault = slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 1})))
}

// Unquiet restores the logger saved by Quiet.
func Unquiet() {
	if savedDefault != nil {
		slog.SetDefault(savedDefault)
		savedDefault = nil
	}
}

func Default() *slog.Logger {
	return slog.Default()
}
