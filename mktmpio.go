package main

import (
	"github.com/codegangsta/cli"
	"github.com/mktmpio/cli/commands"
	"os"
)

var version = "HEAD"

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = "mktmpio"
	app.Usage = "create, destroy, and manage mktmpio instances"
	// Make the default action a remote shell, so 'mktmpio redis' always gives a
	// remote Redis shell
	app.Action = commands.ShellCommand.Action
	app.Commands = []cli.Command{
		commands.ShellCommand,
	}
	app.Run(os.Args)
}
