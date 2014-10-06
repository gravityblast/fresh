package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func checkPath(ignore []string, path string, info os.FileInfo, err error) error {
	if info.IsDir() && !isTmpDir(path) {
		if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}
		for _, ig := range ignore {
			if strings.Contains(path, ig) {
				return filepath.SkipDir
			}
		}
		watchFolder(path)
	}
	return err
}

func watch() {
	root := root()
	ignore := ignoreList()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return checkPath(ignore, path, info, err)
	})
	for _, w := range watchPaths() {
		filepath.Walk(w, func(path string, info os.FileInfo, err error) error {
			return checkPath(ignore, path, info, err)
		})
	}
}
