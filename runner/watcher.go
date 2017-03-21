package runner

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/radovskyb/watcher"
)

type Watcher interface {
	Add(path string) error
}

func fsWatcher() Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Errors:
				watcherLog("error: %s", err)
			case <-stopChannel:
				watcher.Close()
				return
			}
		}
	}()
	return watcher
}

func pollWatcher(d time.Duration) Watcher {
	watcher := watcher.New()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Path) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			case <-stopChannel:
				watcher.Close()
			case <-watcher.Closed:
				return
			}
		}
	}()

	go func() {
		if err := watcher.Start(d); err != nil {
			fatal(err)
		}
	}()
	watcher.Wait()
	return watcher
}

func watch() {
	var watcher Watcher

	poll, duration := pollDuration()
	if poll {
		watcherLog("polling for changes, duration: %s", duration)
		watcher = pollWatcher(duration)
	} else {
		watcher = fsWatcher()
	}

	root := root()
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnoredFolder(path) {
				watcherLog("Ignoring %s", path)
				return filepath.SkipDir
			}

			watcherLog("Watching %s", path)
			if err := watcher.Add(path); err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		fatal(err)
	}
}
