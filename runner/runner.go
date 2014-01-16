package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath())

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

	go func() {
		io.Copy(appLogWriter{}, stderr)
	}()

	go func() {
		io.Copy(appLogWriter{}, stdout)
	}()

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
