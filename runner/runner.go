package runner

import (
	"io"
	"os"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath())

	if len(os.Args) > 2 && os.Args[1] == "-c" {
		cmd.Args = append([]string{buildPath()}, os.Args[3:]...)
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

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
