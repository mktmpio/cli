// Copyright Datajin Technologies, Inc. 2015,2016. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path"
	"text/template"
)

const shortInstanceListTemplate = `{{range .}}{{.ID}}
{{end}}`

const longInstanceListTemplate = `Instances: {{len .}}
{{w "ID"}}{{w "Type"}}{{w "User"   }}{{w "Password"}}{{w "Host" }}{{n "Port" }}
=======================================================================================
{{range .}}{{w .ID }}{{w .Type }}{{w .Username}}{{w .Password }}{{w .Host  }}{{n .Port  }}
{{end}}`

var helpers = template.FuncMap{
	"w": padColumnWide,
	"n": padColumnNarrow,
}

func padColumnWide(v interface{}) string {
	return fmt.Sprintf("%-16v", v)
}

func padColumnNarrow(v interface{}) string {
	return fmt.Sprintf("%-8v", v)
}

// Definition for the 'mktmpio shell' command
var ListCommand = cli.Command{
	Name:   "ls",
	Usage:  "list and inspect running database servers",
	Action: lsAction,
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "long, l"},
	},
}

// shellAction implements the 'mktmpio shell' command
func lsAction(c *cli.Context) {
	instances, err := client.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing instance: %s\n", err)
		return
	}
	t := template.New("instance").Funcs(helpers)
	if c.Bool("long") {
		t = template.Must(t.Parse(longInstanceListTemplate))
	} else {
		t = template.Must(t.Parse(shortInstanceListTemplate))
	}
	selected := instances[0:0:0]
	if len(c.Args()) == 0 {
		selected = append(selected, instances...)
	} else {
		for _, i := range instances {
			for _, arg := range c.Args() {
				if matched, _ := path.Match(arg, i.ID); matched {
					selected = append(selected, i)
				}
			}
		}
	}
	t.Execute(c.App.Writer, selected)
}
