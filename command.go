package main

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

type Command struct {
	Name      string
	CmdString string
	Cmd       *exec.Cmd
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
	Logger    *Logger
}

func newCommand(name, command string) *Command {
	c := &Command{
		Name:      name,
		CmdString: command,
		Logger:    newLogger(name),
	}

	return c
}

func (c *Command) build() error {
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

func (c *Command) Run(done chan bool) {
	logger.log("running command %v\n", c.Name)

	err := c.build()
	if err != nil {
		log.Fatal(err)
	}

	go io.Copy(c.Logger, c.Stdout)
	go io.Copy(c.Logger, c.Stderr)

	err = c.Cmd.Run()
	logger.log("Errors on %s: %v\n", c.CmdString, err)
	done <- true
}

func (c *Command) Stop() {
	if c.Cmd != nil && c.Cmd.Process != nil {
		logger.log("killing process %v\n", c.CmdString)
		c.Cmd.Process.Kill()
	}
}
