package logger

import "log"

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (*Logger) Msg(v ...any) {
	log.Println(v...)
}
