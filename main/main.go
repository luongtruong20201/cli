package main

import (
	"cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rm-cli"
	app.Usage = "A simple CLI application"
	app.Version = "1.0.0"
	app.Author = "Your Name"
	app.Email = "your.email@example.com"

	command := cli.Command{
		Name:      "truonglq",
		ShortName: "tlq",
		Action: func(c *cli.Context) {
			fmt.Println("Hello from truonglq command!")
		},
		BashComplete: func(c *cli.Context) {
			fmt.Println("Completion for truonglq command")
		},
	}
	app.Commands = []cli.Command{command}

	app.Run(os.Args)
}
