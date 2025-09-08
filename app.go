package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type App struct {
	Name                 string
	Usage                string
	Version              string
	Commands             []Command
	Flags                []Flag
	EnableBashCompletion bool
	BashComplete         func(context *Context)
	Before               func(context *Context) error
	Action               func(context *Context)
	CommandNotFound      func(context *Context, command string)
	Compiled             time.Time
	Author               string
	Email                string
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
		Name:         os.Args[0],
		Usage:        "A new cli application",
		Version:      "0.0.0",
		BashComplete: DefaultAppComplete,
		Action:       helpCommand.Action,
		Compiled:     compileTime(),
		Author:       "Author",
		Email:        "unknown@email",
	}
}

func (a *App) Run(arguments []string) error {
	if a.Command(helpCommand.Name) == nil {
		a.Commands = append(a.Commands, helpCommand)
	}
	if a.EnableBashCompletion {
		a.appendFlag(BashCompletionFlag)
	}
	a.appendFlag(VersionFlag)
	a.appendFlag(HelpFlag)
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments[1:])
	nerr := normalizeFlags(a.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		context := NewContext(a, set, set)
		ShowAppHelp(context)
		fmt.Println("")
		return nerr
	}
	context := NewContext(a, set, set)
	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		ShowAppHelp(context)
		fmt.Println("")
		return err
	}
	if checkCompletions(context) {
		return nil
	}
	if checkHelp(context) {
		return nil
	}
	if checkVersion(context) {
		return nil
	}
	if a.Before != nil {
		err := a.Before(context)
		if err != nil {
			return err
		}
	}
	args := context.Args()
	if args.Present() {
		name := args.First()
		c := a.Command(name)
		if c != nil {
			return c.Run(context)
		}
	}
	a.Action(context)
	return nil
}

func (a *App) RunAsSubcommand(ctx *Context) error {
	if len(a.Commands) > 0 {
		if a.Command(helpCommand.Name) == nil {
			a.Commands = append(a.Commands, helpCommand)
		}
	}
	if a.EnableBashCompletion {
		a.appendFlag(BashCompletionFlag)
	}
	a.appendFlag(HelpFlag)
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(ctx.Args().Tail())
	nerr := normalizeFlags(a.Flags, set)
	context := NewContext(a, set, ctx.globalSet)
	if nerr != nil {
		fmt.Println(nerr)
		if len(a.Commands) > 0 {
			ShowSubcommandHelp(context)
		} else {
			ShowCommandHelp(ctx, context.Args().First())
		}
		fmt.Println("")
		return nerr
	}
	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		ShowSubcommandHelp(context)
		return err
	}
	if checkCompletions(context) {
		return nil
	}
	if len(a.Commands) > 0 {
		if checkSubcommandHelp(context) {
			return nil
		}
	} else {
		if checkCommandHelp(ctx, context.Args().First()) {
			return nil
		}
	}
	if a.Before != nil {
		err := a.Before(context)
		if err != nil {
			return err
		}
	}
	args := context.Args()
	if args.Present() {
		name := args.First()
		c := a.Command(name)
		if c != nil {
			return c.Run(context)
		}
	}
	if len(a.Commands) > 0 {
		a.Action(context)
	} else {
		a.Action(ctx)
	}
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
