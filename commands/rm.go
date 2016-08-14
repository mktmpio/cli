// Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/urfave/cli"
)

// RemoveCommand defines the 'mktmpio shell' command
var RemoveCommand = cli.Command{
	Name:   "rm",
	Usage:  "shutdown running database servers",
	Action: rmAction,
}

// rmAction implements the 'mktmpio shell' command
func rmAction(c *cli.Context) {
	var removals sync.WaitGroup
	for _, id := range c.Args() {
		removals.Add(1)
		go func(id string) {
			defer removals.Done()
			err := client.Destroy(id)
			if err == nil {
				fmt.Fprintf(os.Stdout, "Shutdown %s\n", id)
			} else {
				fmt.Fprintf(os.Stderr, "Could not shutdown %s: %v\n", id, err)
			}
		}(id)
	}
	removals.Wait()
}
