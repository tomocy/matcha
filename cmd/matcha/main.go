package main

import "github.com/urfave/cli"

func main() {}

func newApp() *app {
	a := &app{
		App: cli.NewApp(),
	}
	a.setUp()

	return a
}

type app struct {
	*cli.App
}

func (a *app) setUp() {
	a.setCommands()
}

func (a *app) setCommands() {
	a.Commands = []cli.Command{}
}
