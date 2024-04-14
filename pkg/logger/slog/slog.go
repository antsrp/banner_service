package slog

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/antsrp/banner_service/pkg/logger"
)

type Logger struct {
	logger   *slog.Logger
	opts     *slog.HandlerOptions
	isFormat bool
}

var _ logger.Logger = Logger{}

func setupOptions(setups ...func(l *Logger)) Logger {
	l := Logger{
		opts: &slog.HandlerOptions{},
	}

	for _, setup := range setups {
		setup(&l)
	}

	return l
}

func NewTextLogger(w io.Writer, setups ...func(l *Logger)) Logger {
	l := setupOptions(setups...)
	l.logger = slog.New(slog.NewTextHandler(w, l.opts))
	return l
}

func NewJsonLogger(w io.Writer, setups ...func(l *Logger)) Logger {
	l := setupOptions(setups...)
	l.logger = slog.New(slog.NewJSONHandler(w, l.opts))
	return l
}

func (l Logger) prep(msg *string, args *[]any) {
	if l.isFormat {
		*msg = fmt.Sprintf(*msg, *args...)
		*args = nil
	}
}

func (l Logger) Info(msg string, args ...any) {
	l.prep(&msg, &args)
	l.logger.Info(msg, args...)
}
func (l Logger) Error(msg string, args ...any) {
	l.prep(&msg, &args)
	l.logger.Error(msg, args...)
}
func (l Logger) Fatal(msg string, args ...any) {
	l.prep(&msg, &args)
	l.logger.Error(msg, args...)
}
func (l Logger) Debug(msg string, args ...any) {
	l.prep(&msg, &args)
	l.logger.Debug(msg, args...)
}
func (l Logger) Warn(msg string, args ...any) {
	l.prep(&msg, &args)
	l.logger.Warn(msg, args...)
}
