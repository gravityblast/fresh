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
	// mainSection := sections[0]
	// assert.Equal(t, 2, len(mainSection.Commands))
	// assert.Equal(t, "1", mainSection["foo"])
	// assert.Equal(t, "2", mainSection["bar"])

	// // Section 1
	// section1 := sections["section_1"]
	// assert.Equal(t, 5, len(section1))
	// assert.Equal(t, "3", section1["foo"])
	// assert.Equal(t, "4", section1["bar"])
	// assert.Equal(t, "5 6", section1["baz"])
	// assert.Equal(t, "7", section1["qux"])
	// assert.Equal(t, "", section1["quux"])

	// // Section 2
	// section2 := sections["section_2"]
	// assert.Equal(t, 9, len(section2))
	// assert.Equal(t, "1", section2["a"])
	// assert.Equal(t, "2", section2["b"])
	// assert.Equal(t, "3", section2["c"])
	// assert.Equal(t, "4", section2["d"])
	// assert.Equal(t, "5", section2["e"])
	// assert.Equal(t, "6", section2["f"])
	// assert.Equal(t, "7", section2["g"])
	// assert.Equal(t, "8", section2["h"])
	// assert.Equal(t, "http://example.com", section2["url"])

	// // Section 3
	// section3 := sections["section_3"]
	// assert.Equal(t, 0, len(section3))
}
