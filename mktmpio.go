package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio/mktmpio"
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
		err, client := mktmpio.NewClient()
		err, instance := client.Create(c.Args()[0])
		if err != nil {
			fmt.Printf("Error creating instance: %s\n", err)
			return
		}
		defer func() {
			instance.Destroy()
			fmt.Printf("Instance %s terminated.\n", instance.Id)
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
