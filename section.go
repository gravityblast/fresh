package main

import (
	"regexp"
	"strings"
)

type section struct {
	Name     string
	Globs    []string
	Commands []*command
}

func newSection(description string) *section {
	var globsString string
	name := description
	parts := strings.SplitN(description, ":", 2)
	if len(parts) > 1 {
		name = parts[0]
		globsString = parts[1]
	}

	extRe := regexp.MustCompile(`\s*,\s*`)
	globs := []string{}
	for _, rawGlob := range extRe.Split(globsString, -1) {
		glob := strings.TrimSpace(rawGlob)
		if len(glob) > 0 {
			globs = append(globs, glob)
		}
	}

	return &section{
		Name:  name,
		Globs: globs,
	}
}

func (s *section) NewCommand(cmd string) *command {
	c := newCommand(s, cmd)
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
