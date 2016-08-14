// Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mktmpio/go-mktmpio"
	"github.com/urfave/cli"
)

// Config stores the shared mktmpio config used by all the cli commands
var Config = mktmpio.LoadConfig()

var (
	client    *mktmpio.Client
	clientErr error
	logger    = log.New(ioutil.Discard, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
)

// PopulateConfig populates the shared config used by all the cli commands.
func PopulateConfig(c *cli.Context) error {
	if c.GlobalBool("debug") {
		logger.SetOutput(c.App.Writer)
	}
	if c.GlobalIsSet("token") {
		Config.Token = c.GlobalString("token")
	}
	if c.GlobalIsSet("url") {
		Config.URL = c.GlobalString("url")
	}
	logger.Printf("loaded config: %v", Config)
	client, clientErr = mktmpio.NewClient(Config)
	if clientErr != nil {
		fmt.Fprintf(c.App.Writer, "Error initializing client: %s\n", clientErr)
	} else {
		client.UserAgent = "mktmpio-cli/" + c.App.Version + " (go-mktmpio)"
		client.SetLogger(logger)
	}
	logger.Printf("Initialized: %+v", client)
	return clientErr
}
