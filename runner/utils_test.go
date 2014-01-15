package runner

import (
	assert "github.com/pilu/miniassert"
	"testing"
)

func TestIsWatchedFile(t *testing.T) {
	// valid extensions
	assert.True(t, isWatchedFile("test.go"))
	assert.True(t, isWatchedFile("test.tpl"))
	assert.True(t, isWatchedFile("test.tmpl"))
	assert.True(t, isWatchedFile("test.html"))

	/* // invalid extensions */
	assert.False(t, isWatchedFile("test.css"))
	assert.False(t, isWatchedFile("test-executable"))

	// files in tmp
	assert.False(t, isWatchedFile("./tmp/test.go"))
}
