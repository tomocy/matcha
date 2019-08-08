package main

import (
	"fmt"
	"os"

	"github.com/tomocy/matcha/cmd/matcha/client"
	"github.com/urfave/cli"
)

func main() {
	a := newApp()
	if err := a.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %s\n", err)
		os.Exit(1)
	}
}

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
	a.Commands = []cli.Command{
		{
			Name:   "cli",
			Action: a.runCLI,
		},
	}
}

func (a *app) runCLI(ctx *cli.Context) error {
	var c client.CLI
	return c.FetchPosts()
}
