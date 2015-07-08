package main

import (
	"bufio"
	"strings"
	"testing"

	assert "github.com/pilu/miniassert"
)

func TestParseConfig(t *testing.T) {
	content := `
  # comment
  [section 1] # aslid las dlkj s
	WATCH ./public/js
	RUN
	RUN
  [section 2]
  [section 3]
	#WATCH .`

	reader := bufio.NewReader(strings.NewReader(content))

	cs := newConfigScanner(reader)

	config := &config{}
	err := cs.scan(config)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(config.sections))

	s := config.sections[0]
	assert.Equal(t, "MAIN", s.Name)

	s = config.sections[1]
	assert.Equal(t, "section 1", s.Name)
	assert.Equal(t, 2, len(s.Commands))
	assert.Equal(t, "compile-js", s.Commands[0].Name)
	// assert.Equal(t, "compile-js -w ./public/javascripts", s.Commands[0].CmdString)
	// assert.Equal(t, "minify-js", s.Commands[1].Name)
	// assert.Equal(t, "minify-js ./public/javascripts/app.js", s.Commands[1].CmdString)

	s = config.sections[2]
	assert.Equal(t, "section 2", s.Name)

	s = config.sections[3]
	assert.Equal(t, "section 3", s.Name)
}
