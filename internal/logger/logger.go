package logger

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dpotapov/slogpfx"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"

	"github.com/Dydhzo/stremthru-proxy/internal/config"
)

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var _ = func() *slog.Logger {
	w := os.Stderr

	var handler slog.Handler

	if config.LogFormat == "json" {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: parseLogLevel(config.LogLevel),
		})
	} else {
		handler = slogpfx.NewHandler(
			tint.NewHandler(w, &tint.Options{
				Level:      parseLogLevel(config.LogLevel),
				NoColor:    !isatty.IsTerminal(w.Fd()),
				TimeFormat: time.DateTime,
			}),
			&slogpfx.HandlerOptions{
				PrefixKeys: []string{"scope"},
			},
		)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelInfo)
	return logger
}()

func Scoped(scope string) *slog.Logger {
	return slog.With("scope", scope)
}
