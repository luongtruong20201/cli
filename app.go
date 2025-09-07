package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type App struct {
	Name     string
	Usage    string
	Version  string
	Commands []Command
	Flags    []Flag
	Action   func(context *Context)
	Compiled time.Time
	Author   string
	Email    string
}

func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

func NewApp() *App {
	return &App{
		Name:     os.Args[0],
		Usage:    "A new cli application",
		Version:  "0.0.0",
		Action:   helpCommand.Action,
		Compiled: compileTime(),
		Author:   "Author",
		Email:    "author@gmail.com",
	}
}

func (a *App) Run(arguments []string) error {
	if a.Command(helpCommand.Name) == nil {
		a.Commands = append(a.Commands, helpCommand)
	}
	a.appendFlag(BoolFlag{"version", "print the version"})
	a.appendFlag(helpFlag{"show help"})
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments[1:])
	if err := normalizeFlags(a.Flags, set); err != nil {
		fmt.Println(err)
		context := NewContext(a, set, set)
		ShowAppHelp(context)
		fmt.Println("")
		return err
	}
	context := NewContext(a, set, set)

	if err != nil {
		fmt.Println("Incorrect Usage.")
		fmt.Println("")
		ShowAppHelp(context)
		fmt.Println("")
		return err
	}
	if checkHelp(context) {
		return nil
	}
	if checkVersion(context) {
		return nil
	}
	args := context.Args()
	if len(args) > 0 {
		name := args[0]
		c := a.Command(name)
		if c != nil {
			return c.Run(context)
		}
	}
	a.Action(context)
	return nil
}

func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return &c
		}
	}
	return nil
}

func (a *App) hasFlag(flag Flag) bool {
	for _, f := range a.Flags {
		if flag == f {
			return true
		}
	}
	return false
}

func (a *App) appendFlag(flag Flag) {
	if !a.hasFlag(flag) {
		a.Flags = append(a.Flags, flag)
	}
}
