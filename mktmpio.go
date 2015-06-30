package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
	"golang.org/x/crypto/ssh/terminal"
	"io"
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
			fmt.Printf("Error initializing client: %s\n", err)
			return
		}
		instance, err := client.Create(c.Args()[0])
		if err != nil {
			fmt.Printf("Error creating %s instance: %s\n", c.Args()[0], err)
			return
		}
		defer func() {
			if err := instance.Destroy(); err != nil {
				fmt.Printf("Error terminating %s instance %s: %v\n", instance.Type, instance.ID, err)
			} else {
				fmt.Printf("Instance %s terminated.\n", instance.ID)
			}
		}()
		if len(instance.ContainerShell) > 0 {
			if err = remoteShell(client, instance); err != nil {
				fmt.Printf("Error running remote %s shell for %s: %v\n", instance.Type, instance.ID, err)
			}
		} else {
			localShell(instance)
		}
	}

	app.Run(os.Args)
}

func localShell(instance *mktmpio.Instance) error {
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
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running local %s shell for %s: %v\n", instance.Type, instance.ID, err)
		return err
	}
	return nil
}

func remoteShell(client *mktmpio.Client, instance *mktmpio.Instance) error {
	reader, writer, err := client.Attach(instance.ID)
	errs := make(chan error)
	pipe := func(r io.Reader, w io.Writer) {
		buf := make([]byte, 128)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				w.Write(buf[:n])
			}
			if err != nil {
				errs <- err
				return
			}
		}
	}
	if err != nil {
		return err
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)
	go pipe(os.Stdin, writer)
	go pipe(reader, os.Stdout)
	err = <-errs
	if err != io.EOF {
		return err
	}
	return nil
}
