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
	Action          func(context *Context)
	Flags           []Flag
	SkipFlagParsing bool
}

func (c Command) Run(ctx *Context) error {
	c.Flags = append(
		c.Flags,
		BoolFlag{"help, h", "show help"},
	)
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
