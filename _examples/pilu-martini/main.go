package main

import (
	"github.com/pilu/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world - Martini"
	})
	m.Run()
}
