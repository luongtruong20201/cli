package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Command struct {
	Name            string
	ShortName       string
	Usage           string
	Description     string
	BashComplete    func(context *Context)
	Before          func(context *Context) error
	Action          func(context *Context) error
	Subcommands     []Command
	Flags           []Flag
	SkipFlagParsing bool
	HideHelp        bool
}

func (c Command) Run(ctx *Context) error {
	if len(c.Subcommands) > 0 || c.Before != nil {
		return c.startApp(ctx)
	}
	if !c.HideHelp {
		c.Flags = append(
			c.Flags,
			HelpFlag,
		)
	}
	if ctx.App.EnableBashCompletion {
		c.Flags = append(c.Flags, BashCompletionFlag)
	}
	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	firstFlagIndex := -1
	for index, arg := range ctx.Args() {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}
	var err error
	if firstFlagIndex > -1 && !c.SkipFlagParsing {
		args := ctx.Args()
		regularArgs := args[1:firstFlagIndex]
		flagArgs := args[firstFlagIndex:]
		err = set.Parse(append(flagArgs, regularArgs...))
	} else {
		err = set.Parse(ctx.Args().Tail())
	}
	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return err
	}
	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		fmt.Println("")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return nerr
	}
	context := NewContext(ctx.App, set, ctx.globalSet)
	if checkCommandCompletions(context, c.Name) {
		return nil
	}
	if checkCommandHelp(context, c.Name) {
		return nil
	}
	context.Command = c
	c.Action(context)
	return nil
}

func (c Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}

func (c Command) startApp(ctx *Context) error {
	app := NewApp()
	app.Name = fmt.Sprintf("%s %s", ctx.App.Name, c.Name)
	if c.Description != "" {
		app.Usage = c.Description
	} else {
		app.Usage = c.Usage
	}
	app.CommandNotFound = ctx.App.CommandNotFound
	app.Commands = c.Subcommands
	app.Flags = c.Flags
	app.HideHelp = c.HideHelp
	app.EnableBashCompletion = ctx.App.EnableBashCompletion
	if c.BashComplete != nil {
		app.BashComplete = c.BashComplete
	}
	app.Before = c.Before
	if c.Action != nil {
		app.Action = c.Action
	} else {
		app.Action = helpSubcommand.Action
	}
	return app.RunAsSubcommand(ctx)
}
