package logger

import (
	"log"
	"os"
	"sync"
)

type LoggerStd struct {
	log     *log.Logger
	out     *os.File
	outLock *sync.Mutex
}

func NewLoggerStd() *LoggerStd {
	return &LoggerStd{
		log:     log.New(os.Stderr, "", log.LstdFlags),
		out:     nil,
		outLock: &sync.Mutex{},
	}
}

func (l *LoggerStd) Msg(v ...any) {
	l.log.Println(v...)
}

func (l *LoggerStd) OpenFile(fileName string) error {
	new, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.outLock.Lock()
	defer l.outLock.Unlock()
	now := l.out
	l.out = new
	l.log.SetOutput(new)
	if now != nil {
		now.Close()
	}
	return nil
}
