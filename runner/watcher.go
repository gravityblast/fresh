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

func watch() {
	root := root()
	ignorePathsArr := strings.Split(settings["ignore_dirs"], ",")
	ignorePaths := make(map[string]bool)
	for _, path := range ignorePathsArr {
		if strings.HasPrefix(path, " ") || strings.HasSuffix(path, " ") {
			path = strings.Trim(path, " ")
		}
		ignorePaths[path] = true
	}
	watcherLog("Ignore Paths: %v", ignorePaths)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if _, ignore := ignorePaths[info.Name()]; (len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")) || ignore {
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
