package main

import (
	"testing"

	assert "github.com/pilu/miniassert"
)

func TestNewRunner(t *testing.T) {
	r := newRunner()
	assert.Equal(t, 0, len(r.Sections))
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
