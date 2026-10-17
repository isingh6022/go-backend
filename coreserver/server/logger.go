package server

import (
	"log/slog"
	"os"
)

type logLevel struct {
	value string
}

func (l logLevel) String() string {
	return l.value
}

var (
	LevelDebug = logLevel{value: "DEBUG"}
	LevelInfo  = logLevel{value: "INFO"}
	LevelWarn  = logLevel{value: "WARN"}
	LevelError = logLevel{value: "ERROR"}
)

func (l logLevel) toSlogLevel() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func fromSlogLevel(sl slog.Level) logLevel {
	switch sl {
	case slog.LevelDebug:
		return LevelDebug
	case slog.LevelWarn:
		return LevelWarn
	case slog.LevelError:
		return LevelError
	default:
		return LevelInfo
	}
}

type logger struct {
	*slog.Logger
	level *slog.LevelVar
}

func (l *logger) SetLevel(level logLevel) {
	l.level.Set(level.toSlogLevel())
}

func (l *logger) Level() logLevel {
	return fromSlogLevel(l.level.Level())
}

func (l *logger) isDebugEnabled() bool {
	return l.level.Level() <= slog.LevelDebug
}

func CreateLogger(level logLevel) *logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(level.toSlogLevel())

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelVar,
	})

	return &logger{
		Logger: slog.New(jsonHandler),
		level:  levelVar,
	}
}
