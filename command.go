package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type command struct {
	Section   *section
	Name      string
	CmdString string
	Cmd       *exec.Cmd
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
	Logger    *customLogger
}

func newCommand(section *section, name, cmd string) *command {
	loggerPrefix := fmt.Sprintf("%s - %s", section.Name, name)
	c := &command{
		Section:   section,
		Name:      name,
		CmdString: cmd,
		Logger:    newLogger(loggerPrefix),
	}

	return c
}

func (c *command) build() error {
	options := strings.Split(c.CmdString, " ")
	c.Cmd = exec.Command(options[0], options[1:]...)

	var err error
	c.Stdout, err = c.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	c.Stderr, err = c.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	return nil
}

func (c *command) Run() error {
	logger.log("Running command %s: %s\n", c.Name, c.CmdString)

	err := c.build()
	if err != nil {
		log.Fatal(err)
	}

	go io.Copy(c.Logger, c.Stdout)
	go io.Copy(c.Logger, c.Stderr)

	err = c.Cmd.Start()
	if err != nil {
		logger.log("Errors on `%s - %s`: %v\n", c.Section.Name, c.Name, err)
	}

	logger.log(fmt.Sprintf("`%s - %s` started with pid %d", c.Section.Name, c.Name, c.Cmd.Process.Pid))

	err = c.Cmd.Wait()
	if err != nil {
		logger.log("Errors on `%s - %s`: %v\n", c.Section.Name, c.Name, err)
	}

	logger.log("`%s - %s` ended\n", c.Section.Name, c.Name)

	return err
}

func (c *command) Stop() {
	if c.Cmd != nil && c.Cmd.Process != nil {
		logger.log("Killing process `%s`\n", c.Name)
		c.Cmd.Process.Kill()
	}
}
