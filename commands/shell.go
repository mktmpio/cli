// Copyright 2015 Datajin Technologies, Inc. All rights reserved.

package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"time"
)

// Definition for the 'mktmpio shell' command
var ShellCommand = cli.Command{
	Name:   "shell",
	Usage:  "create a new server and attach a shell session to it",
	Action: shellAction,
}

// shellAction implements the 'mktmpio shell' command
func shellAction(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowAppHelp(c)
		return
	}
	cfg := mktmpio.LoadConfig()
	client, err := mktmpio.NewClient(cfg)
	if err != nil {
		fmt.Printf("Error initializing client: %s\n", err)
		return
	}
	client.UserAgent = fmt.Sprintf("mktmpio-cli/%s (go-mktmpio)", c.App.Version)
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

func localShell(instance *mktmpio.Instance) error {
	_ = instance.LoadEnv()
	cmd := instance.Cmd()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// some server types are slower than others
	if t := instance.Type; t == "mysql" || t == "couchdb" {
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

func pipe(r io.Reader, w io.Writer, c chan<- error) {
	_, err := io.Copy(w, r)
	c <- err
}

func remoteShell(client *mktmpio.Client, instance *mktmpio.Instance) error {
	conn, err := client.Attach(instance.ID)
	if err != nil {
		return err
	}
	// Only do raw TTY mode if all of stdio is attached to a TTY
	if terminal.IsTerminal(0) && terminal.IsTerminal(1) && terminal.IsTerminal(2) {
		if oldState, err := terminal.MakeRaw(0); err != nil {
			panic(err)
		} else {
			defer terminal.Restore(0, oldState)
		}
	}
	errs := make(chan error)
	go pipe(os.Stdin, conn, errs)
	go pipe(conn, os.Stdout, errs)
	return <-errs
}
