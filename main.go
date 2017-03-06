/*
Fresh is a command line tool that builds and (re)starts your web application everytime you save a go or template file.

If the web framework you are using supports the Fresh runner, it will show build errors on your browser.

It currently works with Traffic (https://github.com/pilu/traffic), Martini (https://github.com/codegangsta/martini) and gocraft/web (https://github.com/gocraft/web).

Fresh will watch for file events, and every time you create/modifiy/delete a file it will build and restart the application.
If `go build` returns an error, it will logs it in the tmp folder.

Traffic (https://github.com/pilu/traffic) already has a middleware that shows the content of that file if it is present. This middleware is automatically added if you run a Traffic web app in dev mode with Fresh.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/pilu/fresh/runner"
	"os"
)

func main() {
	configPath := flag.String("c", "", "config file path")
	flag.Parse()

	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != nil {
			fmt.Printf("Can't find config file `%s`\n", *configPath)
			os.Exit(1)
		} else {
			os.Setenv("RUNNER_CONFIG_PATH", *configPath)
		}
	}

	runner.Start()
}
