package main

import (
	"regexp"
	"strings"
)

type section struct {
	Name       string
	Extensions []string
	Commands   []*command
}

func newSection(description string) *section {
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

	return &section{
		Name:       name,
		Extensions: extensions,
	}
}

func (s *section) NewCommand(name, cmd string) *command {
	c := newCommand(name, cmd)
	s.Commands = append(s.Commands, c)
	return c
}

func (s *section) Run() {
	logger.log("Running section `%s`", s.Name)
	for _, c := range s.Commands {
		err := c.Run()
		if err != nil {
			break
		}
	}
	logger.log("Section ended `%s`", s.Name)
}

func (s *section) Stop() {
	logger.log("Stopping section %s", s.Name)
	for _, c := range s.Commands {
		c.Stop()
	}
}
