package logger

import "log"

type Logger interface {
	Info(msg string, args ...any)
	Err(msg string, args ...any)
	Debug(msg string, args ...any)
	Fatal(msg string, args ...any)
	PrefixPairs(kv ...any)
	SetPrefix(prefix string)
	WithPrefixPairs(kv ...any)
	WithPrefix(prefix string)
}

type tlogger struct {
	logger *log.Logger
}
