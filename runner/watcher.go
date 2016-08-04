package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string) {
	if !strings.Contains(path, "vendor/") && !strings.Contains(path, "bower_components") && !strings.Contains(path, "node_modules") && !strings.Contains(path, "public/system/") && !strings.Contains(path, ".tmpl") {
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
}

func watch() {
	// root := root()
	paths := root()
	if extraDirs := extraDirs(); extraDirs != "" {
		paths += "," + extraDirs
	}

	for _, root := range strings.Split(paths, ",") {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && !isTmpDir(path) {
				if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}

				watchFolder(path)
			}

			return err
		})
	}
}
