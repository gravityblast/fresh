package main

import (
	assert "github.com/pilu/miniassert"
	"testing"
)

func TestNewRunner(t *testing.T) {
	r := newRunner()
	assert.Equal(t, 0, len(r.Sections))
}

func TestRunner_NewSection(t *testing.T) {
	r := newRunner()
	s := r.NewSection(".go")
	assert.Equal(t, 1, len(r.Sections))
	assert.Equal(t, s, r.Sections[0])
}

func TestNewRunnerWithFreshfile_WithFileNotFound(t *testing.T) {
	_, err := newRunnerWithFreshfile("./_test_fixtures/file-not-found")
	assert.NotNil(t, err)
}

func TestNewRunnerWithFreshfile_WithValidFile(t *testing.T) {
	r, err := newRunnerWithFreshfile("./_test_fixtures/Freshfile")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(r.Sections))
}
