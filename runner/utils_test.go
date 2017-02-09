package runner

import (
	"testing"
)

func TestIsWatchedFile(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.tpl", true},
		{"test.tmpl", true},
		{"test.html", true},
		{"test.css", false},
		{"test-executable", false},
		{"./tmp/test.go", false},
	}

	for _, test := range tests {
		actual := isWatchedFile(test.file)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestShouldRebuild(t *testing.T) {
	tests := []struct {
		eventName string
		expected  bool
	}{
		{"test.go", true},
		{"test.tpl", false},
		{"test.tmpl", false},
		{"test.html", false},
		{"test.css", false},
		{"test-executable", false},
		{"./tmp/test.go", true},
	}

	for _, test := range tests {
		actual := shouldRebuild(test.eventName)

		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestIsIgnoredFolder(t *testing.T) {
	tests := []struct {
		dir      string
		expected bool
	}{
		{"assets/node_modules", true},
		{"tmp/pid", true},
		{"app/controllers", false},
	}

	for _, test := range tests {
		actual := isIgnoredFolder(test.dir)
		if actual != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, actual)
		}
	}
}
