package cli

import (
	"fmt"
	"io/ioutil"
	"os"
)

type App struct {
	Name     string
	Usage    string
	Version  string
	Commands []Command
	Flags    []Flag
	Action   func(context *Context)
}

func NewApp() *App {
	return &App{
		Name:    os.Args[0],
		Usage:   "A new cli application",
		Version: "0.0.0",
		Action:  helpCommand.Action,
	}
}

func (a *App) Run(arguments []string) {
	a.Commands = append(a.Commands, helpCommand)
	a.Flags = append(
		a.Flags,
		BoolFlag{"version", "print the version"},
		helpFlag{"show help"},
	)
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments[1:])
	context := NewContext(a, set, set)

	if err != nil {
		fmt.Println("Incorrect Usage.")
		ShowAppHelp(context)
		fmt.Println("")
		os.Exit(1)
	}
	checkHelp(context)
	checkVersion(context)
	args := context.Args()
	if len(args) > 0 {
		name := args[0]
		for _, c := range a.Commands {
			if c.HasName(name) {
				c.Run(context)
				return
			}
		}
	}
	a.Action(context)
}
