package log

import (
	"fmt"
	"log"
	"os"
)

var (
	logger Logger
)

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func GetLogger() Logger {
	if logger == nil {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		logger = &DefaultLogger{}
	}

	return logger
}

func SetLogger(lg Logger) {
	logger = lg
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Info(v ...interface{}) {
	log.Default().Output(2, fmt.Sprint(v...))
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	log.Default().Output(2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Error(v ...interface{}) {
	log.Default().Output(2, fmt.Sprint(v...))
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	log.Default().Output(2, fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) Fatal(v ...interface{}) {
	log.Default().Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	log.Default().Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

type EmptyLogger struct {
}

func (l *EmptyLogger) Info(v ...interface{}) {
}

func (l *EmptyLogger) Infof(format string, v ...interface{}) {
}

func (l *EmptyLogger) Error(v ...interface{}) {
}

func (l *EmptyLogger) Errorf(format string, v ...interface{}) {
}

func (l *EmptyLogger) Fatal(v ...interface{}) {
	log.Default().Output(2, fmt.Sprint(v...))
	os.Exit(1)

}

func (l *EmptyLogger) Fatalf(format string, v ...interface{}) {
	log.Default().Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
