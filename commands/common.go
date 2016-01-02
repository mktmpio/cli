// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
)

// Config stores the shared mktmpio config used by all the cli commands
var Config = mktmpio.LoadConfig()
var client *mktmpio.Client

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

// InitializeClient returns an error if the shared client is not valid
func InitializeClient(c *cli.Context) error {
	client, err := mktmpio.NewClient(Config)
	if err != nil {
		fmt.Fprintf(c.App.Writer, "Error initializing client: %s\n", err)
	} else {
		client.UserAgent = "mktmpio-cli/" + c.App.Version + " (go-mktmpio)"
	}
	return err
}
