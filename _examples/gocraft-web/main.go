package main

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/pilu/fresh/runner/runnerutils"
	"net/http"
)

func runnerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(rw)
		return
	}

	next(rw, req)
}

type Context struct{}

func (c *Context) SayHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello World")
}

func main() {
	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(runnerMiddleware).
		Get("/", (*Context).SayHello)
	http.ListenAndServe("localhost:3000", router)
}
