package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "mktmpio"
	app.Usage = "create, destroy, and manage mktmpio instances"
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			cli.ShowAppHelp(c)
			return
		}
		client, err := mktmpio.NewClient()
		if err != nil {
			fmt.Printf("Error creating client: %s\n", err)
			return
		}
		instance, err := client.Create(c.Args()[0])
		if err != nil {
			fmt.Printf("Error creating instance: %s\n", err)
			return
		}
		defer func() {
			if err := instance.Destroy(); err != nil {
				fmt.Printf("Error terminatined instance %s: %v\n", instance.ID, err)
			} else {
				fmt.Printf("Instance %s terminated.\n", instance.ID)
			}
		}()
		_ = instance.LoadEnv()
		cmd := instance.Cmd()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// MySQL is particularly slow to start up
		if instance.Type == "mysql" {
			time.Sleep(500 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		err = cmd.Run()
	}

	app.Run(os.Args)
}
