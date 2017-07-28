package runner

import (
	"os"
	"strings"

	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"io/ioutil"
	"os/exec"
	"runtime"
)

var watchedFolders = map[string]struct{}{}
var pkgName = ""
var pkgPath = ""
var goPaths = []string{}
var watcher *fsnotify.Watcher

func initWatcher() {
    // parse running package name
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	goPaths = strings.Split(os.Getenv("GOPATH"), separator)
	root := root()
    dir, err := os.Getwd()
    if err != nil {
        fatal(err)
    }
	if root != "." {
        pkgName = root
        for _, gopath := range goPaths {
            pkgPath = gopath + "/src/" + pkgName
            e, err := exists(pkgPath)
            if err != nil {
                fatal(err)
            }
            if e {
                break
            }
        }
	} else {
        for _, gopath := range goPaths {
            if strings.HasPrefix(dir, gopath) && len(dir) > len(gopath)+5 {
                pkgName = dir[len(gopath)+5:]
                pkgPath = dir
                break
            }
        }
    }

    // init watcher
	watcher, err = fsnotify.NewWatcher()
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
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getFolders() map[string]struct{} {
	cmd := exec.Command("go", "list", "-f", `{{ join .Deps "\n" }}`, pkgName)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	imps, _ := ioutil.ReadAll(stdout)
	io.Copy(os.Stderr, stderr)

	err = cmd.Wait()
	if err != nil {
		fatal(err)
	}

	_imps := strings.Split(strings.Trim(string(imps), "\n"), "\n")
	_watchedFolders := map[string]struct{}{pkgPath: struct{}{}}
	for _, imp := range _imps {
		for _, gopath := range goPaths {
			path := fmt.Sprintf("%s/src/%s", gopath, imp)
			e, err := exists(path)
			if err != nil {
				fatal(err)
			}
			if e {
				_watchedFolders[path] = struct{}{}
				break
			}
		}
	}

	return _watchedFolders
}

func watch() {
	_watchedFolders := getFolders()
	for folder, _ := range watchedFolders {
		_, ok := _watchedFolders[folder]
		if !ok {
			watcherLog("remove watch %s", folder)
			err := watcher.RemoveWatch(folder)
			if err != nil {
				fatal(err)
			}
		}
	}
	for folder, _ := range _watchedFolders {
		_, ok := watchedFolders[folder]
		if !ok && !isTmpDir(folder) {
            if isIgnoredFolder(folder) {
                watcherLog("Ignoring %s", folder)
                continue
            }
			watcherLog("add watch %s", folder)
			err := watcher.Watch(folder)
			if err != nil {
				fatal(err)
			}
		}
	}
	watchedFolders = _watchedFolders
}
