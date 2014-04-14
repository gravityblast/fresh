package main

import (
	"flag"
	"os"
)

var logger *Logger

func init() {
	logger = newLogger("fresh")
	flag.BoolVar(&logger.Verbose, "verbose", false, "verbose")
	flag.Parse()
}

func main() {
	r, err := newRunnerWithFreshfile("Freshfile")
	if err != nil {
		logger.log("%s\n", err.Error())
		os.Exit(1)
	}

	logger.log("Process %d", os.Getpid())
	r.Run()
	<-r.DoneChan
	println("the end")
}
