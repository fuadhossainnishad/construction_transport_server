package logger

import "log"

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...any) {
	log.Println(msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields ...any) {
	log.Println(msg, fields)
}
