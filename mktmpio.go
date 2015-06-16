package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "mktmpio"
  app.Usage = "create, destroy, and manage mktmpio instances"
  app.Action = func(c *cli.Context) {
    println("Coming soon!")
  }

  app.Run(os.Args)
}
