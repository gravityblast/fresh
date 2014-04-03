package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	Name 		string
	Verbose	bool
	*log.Logger
}

func (logger *Logger) Write(p []byte) (n int, err error) {
	logger.log(string(p))
	return len(p), nil
}

var (
	loggersCount        int
	loggerMaxNameLength int
)

func newLogger(name string) *Logger {
	_logger := log.New(os.Stderr, "", 0)

	logger := &Logger{
		Name:   	name,
		Logger: 	_logger,
		Verbose: 	true,
	}

	loggersCount++
	if length := len(name); length > loggerMaxNameLength {
		loggerMaxNameLength = length
	}

	return logger
}

func (logger *Logger) log(format string, v ...interface{}) {
	if !logger.Verbose {
		return
	}
	now := time.Now()
	timeString := fmt.Sprintf("%d:%d:%02d", now.Hour(), now.Minute(), now.Second())
	// format = fmt.Sprintf("%s%s %s |%s %s", color, timeString, prefix, clear, format)
	// logger.Logger.Printf(format, v...)
	prefix := fmt.Sprintf("%s %-10s |", timeString, logger.Name)
	format = fmt.Sprintf("%s %s", prefix, format)
	logger.Logger.Printf(format, v...)
}
