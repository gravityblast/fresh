package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/pilu/config"
)

type Runner struct {
	Sections []*Section
	DoneChan chan bool
}

func newRunner() *Runner {
	return &Runner{
		DoneChan: make(chan bool),
	}
}

func newRunnerWithFreshfile(freshfilePath string) (*Runner, error) {
	r := newRunner()

	sections, err := config.ParseFile(freshfilePath, "main: *")
	if err != nil {
		return r, err
	}

	for s, opts := range sections {
		section := r.NewSection(s)
		for name, cmd := range opts {
			section.NewCommand(name, cmd)
		}
	}

	return r, nil
}

func (r *Runner) NewSection(description string) *Section {
	s := newSection(description)
	r.Sections = append(r.Sections, s)
	return s
}

func (r *Runner) Run() {
	logger.log("Running...")
	logger.log("%d goroutines", runtime.NumGoroutine())
	go r.ListenForSignals()

	for _, s := range r.Sections {
		go func(s *Section) {
			s.Run()
		}(s)
	}
}

func (r *Runner) Stop() {
	logger.log("Stopping all sections")
	for _, s := range r.Sections {
		s.Stop()
	}
}

func (r *Runner) ListenForSignals() {
	logger.log("Listening for signals")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	signal.Notify(sc, os.Kill)
	<-sc
	fmt.Printf("Interrupt a second time to quit\n")
	logger.log("Waiting for a second signal")
	select {
	case <-sc:
		logger.log("Second signal received")
		r.DoneChan <- true
	case <-time.After(1 * time.Second):
		logger.log("Timeout")
		logger.log("Stopping...")
		r.Stop()
		logger.log("Calling Run...")
		r.Run()
	}
}
