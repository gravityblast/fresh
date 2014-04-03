package main

import (
	"regexp"
	"strings"
)

type Section struct {
	Name       string
	Extensions []string
	Commands   []*Command
	StopChan	 chan bool
}

func newSection(description string, stopChan chan bool) *Section {
	var name string
	parts := strings.Split(description, ":")
	if len(parts) > 1 {
		name = parts[0]
		description = parts[1]
	}

	extRe := regexp.MustCompile(`\s*,\s*`)
	extensions := []string{}
	for _, rawExtension := range extRe.Split(description, -1) {
		extension := strings.TrimSpace(rawExtension)
		if len(extension) > 0 {
			extensions = append(extensions, extension)
		}
	}

	return &Section{
		Name:       name,
		Extensions: extensions,
		StopChan:		stopChan,
	}
}

func (s *Section) NewCommand(name, command string) *Command {
	c := newCommand(name, command)
	s.Commands = append(s.Commands, c)
	return c
}

func (s *Section) Run() {
	logger.log("Running section %s", s.Name)
	for _, c := range s.Commands {
		done := make(chan bool)
		go func(c *Command) {
			c.Run(done)
		}(c)
		select {
		case <-done:
			continue
		case <-s.StopChan:
			return
		}
	}
}

func (s *Section) Stop() {
	logger.log("Stopping section %s", s.Name)
	for _, c := range s.Commands {
		c.Stop()
	}
}
