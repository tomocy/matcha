package main

import "github.com/urfave/cli"

func main() {}

type app struct {
	*cli.App
}

func (a *app) setUp() {
	a.setCommands()
}

func (a *app) setCommands() {
	a.Commands = []cli.Command{}
}
