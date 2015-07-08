package main

import (
	"bufio"
	"fmt"
	"text/scanner"
)

type config struct {
	sections []*section
}

type stateFunc func(*config, *section) (stateFunc, *section, error)

type configScanner struct {
	s        scanner.Scanner
	state    stateFunc
	commands map[string]stateFunc
}

func newConfigScanner(r *bufio.Reader) *configScanner {
	var sc scanner.Scanner
	sc.Init(r)

	s := &configScanner{s: sc}
	s.init()

	return s
}

func (s *configScanner) init() {
	s.commands = map[string]stateFunc{
		"RUN":   s.scanCMDRun,
		"WATCH": s.scanCMDWatch,
	}
}

func (s *configScanner) next() rune {
	r := s.s.Next()
	if r != '#' {
		return r
	}

	for r != '\n' && r != scanner.EOF {
		r = s.s.Next()
	}

	return r
}

func (s *configScanner) scan(c *config) error {
	var err error

	sec := &section{Name: "MAIN"}
	c.sections = append(c.sections, sec)

	for s.state = s.scanLine; s.state != nil; {
		s.state, sec, err = s.state(c, sec)
		if err != nil {
			break
		}
	}

	return err
}

func (s *configScanner) scanLine(c *config, sec *section) (stateFunc, *section, error) {
	r := s.s.Peek()
	for r != scanner.EOF {
		if r != ' ' && r != '\t' && r != '\n' && r != '#' {
			if r == '[' {
				return s.scanSection, sec, nil
			}

			return s.scanCMD, sec, nil
		}

		s.next()
		r = s.s.Peek()

	}

	return nil, sec, nil
}

func (s *configScanner) scanSection(c *config, sec *section) (stateFunc, *section, error) {
	r := s.next()
	if r != '[' {
		return nil, sec, s.errorExpectedRune("[")
	}

	sec = &section{}
	c.sections = append(c.sections, sec)

	return s.scanSectionName, sec, nil
}

func (s *configScanner) scanSectionName(c *config, sec *section) (stateFunc, *section, error) {
	var name string
	r := s.s.Peek()
	for r != ']' {
		if r == scanner.EOF || r == '\n' || r == '#' {
			return nil, sec, s.errorExpectedRune("]")
		}

		r = s.s.Next()
		name += string(r)
		r = s.s.Peek()
	}
	s.next()
	sec.Name = name

	return s.scanLine, sec, nil
}

func (s *configScanner) scanCMD(c *config, sec *section) (stateFunc, *section, error) {
	var name string
	r := s.next()
	for r != scanner.EOF {
		if r == ' ' || r == '\t' || r == '\n' {
			break
		}

		name += string(r)
		r = s.next()
	}

	if cmd, ok := s.commands[name]; ok {
		return cmd, sec, nil
	}

	return nil, sec, fmt.Errorf("Unknown command `%s`", name)
}

func (s *configScanner) scanCMDRun(c *config, sec *section) (stateFunc, *section, error) {
	var cmdString string
	r := s.next()
	for r != scanner.EOF && r != '\n' {
		cmdString = cmdString + string(r)
		r = s.next()
	}

	sec.NewCommand(cmdString)

	return s.scanLine, sec, nil
}

func (s *configScanner) scanCMDWatch(c *config, sec *section) (stateFunc, *section, error) {
	var cmd string
	r := s.next()
	for r != scanner.EOF && r != '\n' {
		cmd = cmd + string(r)
		r = s.next()
	}

	return s.scanLine, sec, nil
}

func (s *configScanner) errorExpectedRune(c string) error {
	p := s.s.Pos()
	return fmt.Errorf("Expected `%s` at line %d, col %d", c, p.Line, p.Column)
}
