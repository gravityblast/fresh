package runner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func listen() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTSTP)

	go func() {
		for range c {
			fmt.Println("")
			listenerLog("Force rebuild")
			startChannel <- "force"
		}
	}()

	listenerLog("Listening ctrl+Z")
}
