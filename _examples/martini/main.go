package main

import (
	"github.com/codegangsta/martini"
	"github.com/pilu/fresh/runner/runnerutils"
	"net/http"
	"os"
)

func runnerMiddleware(w http.ResponseWriter, r *http.Request) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(w)
	}
}

func main() {
	m := martini.Classic()

	if os.Getenv("DEV_RUNNER") == "1" {
		m.Use(runnerMiddleware)
	}

	m.Get("/", func() string {
		return "Hello world - Martini"
	})
	m.Run()
}
