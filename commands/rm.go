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
	Name:      "rm",
	Usage:     "shutdown running database servers",
	Action:    rmAction,
	Before:    rmCheckArgs,
	ArgsUsage: "<ID...>",
}

func rmCheckArgs(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError("missing arguments", 1)
	}
	return nil
}

// rmAction implements the 'mktmpio shell' command
func rmAction(c *cli.Context) error {
	var removals sync.WaitGroup
	errs := make([]error, 0, len(c.Args()))
	for _, id := range c.Args() {
		removals.Add(1)
		go func(id string) {
			defer removals.Done()
			err := client.Destroy(id)
			if err == nil {
				fmt.Fprintf(os.Stdout, "Shutdown %s\n", id)
			} else {
				errs = append(errs, err)
				fmt.Fprintf(os.Stderr, "Could not shutdown %s: %v\n", id, err)
			}
		}(id)
	}
	removals.Wait()
	if len(errs) > 0 {
		return cli.NewMultiError(errs...)
	}
	return nil
}
