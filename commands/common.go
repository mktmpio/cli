// Copyright 2015 Datajin Technologies, Inc. All rights reserved.

package commands

import (
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
)

// Config stores the shared mktmpio config used by all the cli commands
var Config = mktmpio.LoadConfig()

// PopulateConfig populates the shared config used by all the cli commands.
func PopulateConfig(c *cli.Context) error {
	if c.GlobalIsSet("token") {
		Config.Token = c.GlobalString("token")
	}
	if c.GlobalIsSet("url") {
		Config.URL = c.GlobalString("url")
	}
	return nil
}
