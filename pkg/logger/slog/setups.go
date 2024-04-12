package slog

import "log/slog"

func WithSource() func(Logger) {
	return func(l Logger) {
		l.opts.AddSource = true
	}
}

func WithDebugLevel() func(Logger) {
	return func(l Logger) {
		l.opts.Level = slog.LevelDebug
	}
}

func WithErrorLevel() func(Logger) {
	return func(l Logger) {
		l.opts.Level = slog.LevelError
	}
}

func WithInfoLevel() func(Logger) {
	return func(l Logger) {
		l.opts.Level = slog.LevelInfo
	}
}

func WithWarnLevel() func(Logger) {
	return func(l Logger) {
		l.opts.Level = slog.LevelWarn
	}
}
