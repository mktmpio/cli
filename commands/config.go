// Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"strings"

	"github.com/mktmpio/go-mktmpio"
	"github.com/urfave/cli"
)

// ConfigCommand is the definition for the 'mktmpio config' command
var ConfigCommand = cli.Command{
	Name:  "config",
	Usage: "view and modify your mktmpio config",
	Subcommands: []cli.Command{
		{
			Name:      "get",
			Usage:     "view config values",
			ArgsUsage: "[NAMES...]",
			Action:    getConfigs,
		},
		{
			Name:      "set",
			Usage:     "set config values",
			ArgsUsage: "<NAME> <VALUE>..., to unset a value use an empty string as the value",
			Action:    setConfigs,
			Before:    setCheckArgs,
		},
	},
}

func getConfigs(c *cli.Context) error {
	args := []string(c.Args())
	for _, key := range args {
		switch strings.ToLower(key) {
		case "token":
			println(Config.Token)
		case "url":
			println(Config.URL)
		default:
			println("Unknown config key:", key)
		}
	}
	if len(args) == 0 {
		print(Config.String())
	}
	return nil
}

func setCheckArgs(c *cli.Context) error {
	nargs := c.NArg()
	if nargs == 0 || nargs%2 != 0 {
		return cli.NewExitError("incorrect number of arguments", 1)
	}
	return nil
}

func setConfigs(c *cli.Context) error {
	args := []string(c.Args())
	for i := 0; i < len(args); i += 2 {
		switch strings.ToLower(args[i]) {
		case "token":
			Config.Token = args[i+1]
		case "url":
			Config.URL = args[i+1]
		default:
			println("Ignoring unknown key:", args[i])
		}
	}
	err := Config.Save(mktmpio.ConfigPath())
	if err != nil {
		println("Error saving config:", err)
	}
	return err
}
