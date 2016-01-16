package runner

import (
	"testing"

	assert "github.com/pilu/miniassert"
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

func TestIsIgnoredFolder(t *testing.T) {
	assert.True(t, isIgnoredFolder("assets/node_module"))
	assert.False(t, isIgnoredFolder("app/controllers"))
	assert.True(t, isIgnoredFolder("tmp/pid"))
}
