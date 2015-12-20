// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
	"strings"
)

// Definition for the 'mktmpio config' command
var ConfigCommand = cli.Command{
	Name:  "config",
	Usage: "view and modify your mktmpio config",
	Subcommands: []cli.Command{
		cli.Command{
			Name:      "get",
			Usage:     "view config values",
			ArgsUsage: "[NAMES...]",
			Action:    getConfigs,
		},
		cli.Command{
			Name:      "set",
			Usage:     "set config values",
			ArgsUsage: "<NAME> <VALUE>..., to unset a value use an empty string as the value",
			Action:    setConfigs,
		},
	},
}

func getConfigs(c *cli.Context) {
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
}

func setConfigs(c *cli.Context) {
	args := []string(c.Args())
	if len(args) == 0 || len(args)%2 != 0 {
		println("invalid number of arguments:", c.Command.ArgsUsage)
		return
	}
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
}
