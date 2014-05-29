package main

import (
	"flag"
	"fmt"
	"os"
)

var logger *Logger

func init() {
	logger = newLoggerWithColor("fresh", "white")
}

func main() {
	var freshfilePath string

	flag.BoolVar(&logger.Verbose, "v", false, "verbose")
	flag.StringVar(&freshfilePath, "f", "./Freshfile", "Freshfile path")
	flag.Parse()

	r, err := newRunnerWithFreshfile(freshfilePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	logger.log("Process %d", os.Getpid())
	r.Run()
	<-r.DoneChan
	println("the end")
}
