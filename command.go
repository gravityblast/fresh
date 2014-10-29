package main

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

type command struct {
	Name      string
	CmdString string
	Cmd       *exec.Cmd
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
	Logger    *customLogger
}

func newCommand(name, cmd string) *command {
	c := &command{
		Name:      name,
		CmdString: cmd,
		Logger:    newLogger(name),
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
	logger.log("Running command %v\n", c.Name)

	err := c.build()
	if err != nil {
		log.Fatal(err)
	}

	go io.Copy(c.Logger, c.Stdout)
	go io.Copy(c.Logger, c.Stderr)

	err = c.Cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (c *command) Stop() {
	if c.Cmd != nil && c.Cmd.Process != nil {
		logger.log("Killing process `%s`\n", c.Name)
		c.Cmd.Process.Kill()
	}
}
