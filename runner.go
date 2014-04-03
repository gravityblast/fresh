package main

import (
	"fmt"
	"github.com/pilu/config"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

type Runner struct {
	Sections []*Section
	DoneChan chan bool
	StopChan chan bool
}

func newRunner() *Runner {
	return &Runner{
		StopChan: make(chan bool),
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
	s := newSection(description, r.StopChan)
	r.Sections = append(r.Sections, s)
	return s
}

func (r *Runner) Run(done chan bool) {
	var wg sync.WaitGroup
	r.Stop()
	logger.log("Running...")
	logger.log("%d goroutines", runtime.NumGoroutine())
	go r.ListenForSignals(done)

	for _, s := range r.Sections {
		wg.Add(1)
		go func(s *Section) {
			defer wg.Done()
			s.Run()
		}(s)
	}
	wg.Wait()
	logger.log("finish")
}

func (r *Runner) Stop() {
	logger.log("stopping all sections")
	for _, s := range r.Sections {
		s.Stop()
	}
}

func (r *Runner) ListenForSignals(done chan bool) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
	fmt.Printf("Interrupt a second time to quit\n")
	select {
	case <-sc:
		r.StopChan <- true
		done <- true
	case <-time.After(1 * time.Second):
		r.StopChan <- true
		r.Stop()
		r.Run(done)
	}
}
