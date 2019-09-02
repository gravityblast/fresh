package runner

import (
	"container/list"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func build() (string, bool) {
	buildLog("Building...")

	args := getBuildArgs()
	cmd := exec.Command("go", args...)

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

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}

func getListOfArgs() (l *list.List) {
	l = list.New()
	build := l.PushBack("build")
	l.PushBack("-o")
	l.PushBack(buildPath())
	l.PushBack(root())

	if settings["race"] == "1" {
		buildLog("Building with --race")
		l.InsertAfter("--race", build)
	}

	return
}

func getBuildArgs() (args []string) {
	list_args := getListOfArgs()

	args = make([]string, list_args.Len())
	var x int = 0
	for e := list_args.Front(); e != nil; e = e.Next() {
		args[x] = e.Value.(string)
		x++
	}

	return
}
