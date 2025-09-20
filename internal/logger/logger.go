package logger

type Logger interface {
	Msg(v ...any)
	OpenFile(fileName string) error
}

func NewLogger() Logger {
	return NewLoggerStd()
}
