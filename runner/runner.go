package runner

import (
  "flag"
  "io"
  "os/exec"
  "strings"
)

var cmdArgs string

func init() {
  flag.StringVar(&cmdArgs, "a", "", "Command line arguments to pass to the process")
}

func run() bool {
  runnerLog("Running...")

  var cmd *exec.Cmd
  if cmdArgs != "" {
    args := strings.Split(cmdArgs, " ")
    cmd = exec.Command(buildPath(), args...)
  } else {
    cmd = exec.Command(buildPath())
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
