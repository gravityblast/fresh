package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

var loggerColors = map[string]int{
	"black":   30,
	"red":     31,
	"green":   32,
	"yellow":  33,
	"blue":    34,
	"magenta": 35,
	"cyan":    36,
	"white":   37,
}

var loggerAvailableColors = []string{
	"cyan",
	"yellow",
	"green",
	"magenta",
	"red",
	"blue",
}

type customLogger struct {
	Name    string
	Verbose bool
	Color   int
	*log.Logger
}

func (logger *customLogger) Write(p []byte) (n int, err error) {
	logger.log(string(p))
	return len(p), nil
}

var (
	loggerColorIndex    int
	loggerMaxNameLength int
)

func newLogger(name string) *customLogger {
	colorIndex := int(math.Mod(float64(loggerColorIndex), float64(len(loggerAvailableColors))))
	colorName := loggerAvailableColors[colorIndex]

	loggerColorIndex++

	if length := len(name); length > loggerMaxNameLength {
		loggerMaxNameLength = length
	}

	return newLoggerWithColor(name, colorName)
}

func newLoggerWithColor(name, colorName string) *customLogger {
	return &customLogger{
		Name:    name,
		Logger:  log.New(os.Stderr, "", 0),
		Verbose: true,
		Color:   loggerColors[colorName],
	}
}

func (logger *customLogger) log(format string, v ...interface{}) {
	if !logger.Verbose {
		return
	}
	now := time.Now()
	timeString := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	formatPadding := fmt.Sprintf("%%-%ds", loggerMaxNameLength)
	prefix := fmt.Sprintf(formatPadding, logger.Name)
	format = fmt.Sprintf("\033[%dm%s %s | \033[0m%s", logger.Color, timeString, prefix, format)
	logger.Logger.Printf(format, v...)
}
