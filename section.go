package main

import (
	"regexp"
	"strings"
)

type Section struct {
	Name       string
	Extensions []string
	Commands   []*Command
}

func newSection(description string) *Section {
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
	}
}

func (s *Section) NewCommand(name, command string) *Command {
	c := newCommand(name, command)
	s.Commands = append(s.Commands, c)
	return c
}

func (s *Section) Run() {
	logger.log("Running section `%s`", s.Name)
	for _, c := range s.Commands {
		err := c.Run()
		if err != nil {
			break
		}
	}
	logger.log("Section ended `%s`", s.Name)
}

func (s *Section) Stop() {
	logger.log("Stopping section %s", s.Name)
	for _, c := range s.Commands {
		c.Stop()
	}
}
