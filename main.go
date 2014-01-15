/*
Fresh is a command line tool that builds and (re)starts your web application everytime you save a go or template file.

If the web framework you are using supports Fresh, it will show build errors on your browser.

It currently works with Traffic (https://github.com/pilu/traffic), Martini (https://github.com/codegangsta/martini) and gocraft/web (https://github.com/gocraft/web).

When you start the `fresh` from your web app folder, Fresh builds the application and (re)starts it. It does the same every time a file is created/modified.
If the build fails, it logs the error message in a file in the tmp folder.
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
