package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var commentSplitRegexp = regexp.MustCompile(`[#;]`)

var keyValueSplitRegexp = regexp.MustCompile(`(\s*(:|=)\s*)|\s+`)

func cleanConfigLine(line string) string {
	chunks := commentSplitRegexp.Split(line, 2)
	return strings.TrimSpace(chunks[0])
}

func parseConfig(reader *bufio.Reader, mainSectionName string) ([]*section, error) {
	var sections []*section
	s := newSection("MAIN")

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return sections, err
		}

		line = cleanConfigLine(line)

		if len(line) == 0 {
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			sections = append(sections, s)
			sectionName := line[1:(len(line) - 1)]
			s = newSection(sectionName)
		} else {
			values := keyValueSplitRegexp.Split(line, 2)
			key := values[0]
			value := ""
			if len(values) == 2 {
				value = values[1]
			}

			s.NewCommand(key, value)
		}
	}

	sections = append(sections, s)

	return sections, nil
}

// func parseConfigFile(path string, mainSectionName string) (configSections, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return make(configSections), err
// 	}

// 	defer file.Close()

// 	reader := bufio.NewReader(file)

// 	return parseConfig(reader, mainSectionName)
// }
