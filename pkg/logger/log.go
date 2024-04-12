package logger

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
	Fatal(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
}
