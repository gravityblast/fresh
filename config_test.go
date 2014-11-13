package main

import (
	"bufio"
	"strings"
	"testing"

	assert "github.com/pilu/miniassert"
)

func TestParseConfig(t *testing.T) {
	content := `
  # comment 1
  ; comment 2

  foo 1
  bar 2

  [section_1]

  foo       3 # using spaces after the key
  bar				4 # using tabs after the key
  # other options for section_1 after section_2

  [section_2]
  a:1
  b: 2
  c : 3
  d :4
  e=5
  f= 6
  g = 7
  h =8

  url: http://example.com

  [section_3]
  `

	reader := bufio.NewReader(strings.NewReader(content))
	sections, err := parseConfig(reader, "main")

	assert.Nil(t, err)
	assert.Equal(t, 4, len(sections))

	// Main section
	mainSection := sections[0]
	assert.Equal(t, 2, len(mainSection.Commands))
	tests := [][]string{
		{"foo", "1"},
		{"bar", "2"},
	}

	for i, opts := range tests {
		assert.Equal(t, opts[0], mainSection.Commands[i].Name)
		assert.Equal(t, opts[1], mainSection.Commands[i].CmdString)
	}

	// Section 1
	section1 := sections[1]
	assert.Equal(t, "section_1", section1.Name)
	tests = [][]string{
		{"foo", "3"},
		{"bar", "4"},
	}

	for i, opts := range tests {
		assert.Equal(t, opts[0], section1.Commands[i].Name)
		assert.Equal(t, opts[1], section1.Commands[i].CmdString)
	}

	// Section 2
	section2 := sections[2]
	assert.Equal(t, "section_2", section2.Name)
	assert.Equal(t, 9, len(section2.Commands))
	tests = [][]string{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
		{"d", "4"},
		{"e", "5"},
		{"f", "6"},
		{"g", "7"},
		{"h", "8"},
		{"url", "http://example.com"},
	}

	for i, opts := range tests {
		assert.Equal(t, opts[0], section2.Commands[i].Name)
		assert.Equal(t, opts[1], section2.Commands[i].CmdString)
	}

	// Section 3
	section3 := sections[3]
	assert.Equal(t, "section_3", section3.Name)
	assert.Equal(t, 0, len(section3.Commands))
}
