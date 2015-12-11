package runner

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func build() (string, bool) {
	var cmd *exec.Cmd
	buildLog("Building...")
	if flag.Lookup("i").Value.String() == "true" {
		buildLog("Installing...")
		cmd = exec.Command("go", "build", "-i", "-o", buildPath(), root())
	} else {
		cmd = exec.Command("go", "build", "-o", buildPath(), root())
	}

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
