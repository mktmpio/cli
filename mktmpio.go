// Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path"
	"runtime"
	"time"

	"github.com/mktmpio/cli/commands"
	"github.com/urfave/cli"
)

// overriden at compile time (-ldflags "-X main.version=V main.commit=C")
var (
	version   = "0.0.0"
	commit    = "HEAD"
	buildtime = "0000-00-00T00:00:00Z"
	t, terr   = time.Parse("2006-01-02T15:04:05Z", buildtime)
)

const appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}}
   {{- if .Flags }} [global options] {{- end -}}
   {{- if .Commands}} command [command options] {{- end -}}
   {{- if .ArgsUsage}} {{.ArgsUsage}} {{- else }} [arguments...] {{- end}}

GLOBAL OPTIONS:
   {{range .Flags}}
     {{- .}}
   {{end}}
COMMANDS:
   {{range .Commands}}
     {{- join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}
BUGS:
   Report to https://github.com/mktmpio/cli/issues

VERSION:
   Version: {{.Version}}
   Compiled: {{.Compiled}}

COPYRIGHT:
   {{.Copyright}}
`

func mktmpioApp() *cli.App {
	// overrides for some variables exposed by urfave/cli
	cli.AppHelpTemplate = appHelpTemplate
	cli.VersionFlag.Name = "version"
	cli.HelpFlag.Name = "help"
	return &cli.App{
		Name:         "mktmpio",
		HelpName:     path.Base(os.Args[0]),
		Usage:        "create, destroy, and manage mktmpio database servers",
		Version:      version + " (built with " + runtime.Compiler + ", " + runtime.Version() + ")",
		Compiled:     t,
		Copyright:    "Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.",
		BashComplete: cli.DefaultAppComplete,
		Action:       commands.ShellCommand.Action,
		Before:       commands.PopulateConfig,
		Writer:       os.Stdout,
		Commands: []cli.Command{
			commands.ConfigCommand,
			commands.ListCommand,
			commands.RemoveCommand,
			commands.ShellCommand,
			commands.LegalCommand,
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "debug, d",
				Usage:  "Enable extra verbose logging",
				EnvVar: "MKTMPIO_DEBUG",
			},
			cli.StringFlag{
				Name:   "token",
				Usage:  "API token for making requests to mktmpio service",
				EnvVar: "MKTMPIO_TOKEN",
				Value:  "TOKEN",
			},
			cli.StringFlag{
				Name:   "url",
				Usage:  "override the URL for the mktmpio service",
				EnvVar: "MKTMPIO_URL",
				Value:  "URL",
			},
		},
		Authors: []cli.Author{
			{Name: "Ryan Graham", Email: "mktmpio@datajin.com"},
		},
	}
}

func main() {
	mktmpioApp().Run(os.Args)
}
