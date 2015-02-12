package runner

import (
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath())

	cmd.Stdout = appLogWriter{}
	cmd.Stderr = appLogWriter{}

	err := cmd.Start()
	if err != nil {
		fatal(err)
	}

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
