// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mktmpio/go-mktmpio"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"sync"
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
	client, err := mktmpio.NewClient(Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing client: %s\n", err)
		return
	}
	client.UserAgent = fmt.Sprintf("mktmpio-cli/%s (go-mktmpio)", c.App.Version)
	instance, err := client.Create(c.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s instance: %s\n", c.Args()[0], err)
		return
	}
	defer func() {
		if err := instance.Destroy(); err != nil {
			fmt.Fprintf(os.Stderr, "Error terminating %s instance %s: %v\n", instance.Type, instance.ID, err)
		} else {
			fmt.Fprintf(os.Stderr, "Instance %s terminated.\n", instance.ID)
		}
	}()
	if len(instance.ContainerShell) > 0 {
		if t := instance.Type; t == "mysql" || t == "couchdb" {
			time.Sleep(500 * time.Millisecond)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		if terminal.IsTerminal(0) && terminal.IsTerminal(1) && terminal.IsTerminal(2) {
			err = remoteShell(client, instance)
		} else {
			err = remoteCmd(client, instance)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running remote %s shell for %s: %v\n", instance.Type, instance.ID, err)
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
		fmt.Fprintf(os.Stderr, "Error running local %s shell for %s: %v\n", instance.Type, instance.ID, err)
		return err
	}
	return nil
}

func pipe(name string, r io.Reader, w io.Writer, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	io.CopyBuffer(w, r, make([]byte, 16))
	// println("pipe().", name)
}

func pipeAndClose(name string, r io.Reader, w io.WriteCloser, wg *sync.WaitGroup) {
	defer w.Close()
	pipe(name, r, w, wg)
	// println("pipeAndClose().", name)
}

func remoteCmd(client *mktmpio.Client, instance *mktmpio.Instance) error {
	stdin, stdout, stderr, err := client.AttachStdio(instance.ID)
	if err != nil {
		return err
	}
	// still not sure why, but this slight delay is needed before writing to the
	// websocket
	time.Sleep(100 * time.Millisecond)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go pipe("stdout", stdout, os.Stdout, &wg)
	go pipe("stderr", stderr, os.Stderr, &wg)
	go pipeAndClose("stdin", os.Stdin, stdin, &wg)
	wg.Wait()
	return nil
}

func remoteShell(client *mktmpio.Client, instance *mktmpio.Instance) error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)
	conn, err := client.Attach(instance.ID)
	if err != nil {
		return err
	}
	go pipe("stdin", os.Stdin, conn, nil)
	pipe("stdout", conn, os.Stdout, nil)
	return conn.Close()
}
