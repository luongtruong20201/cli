package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Command struct {
	Name        string
	ShortName   string
	Usage       string
	Description string
	Action      func(context *Context)
	Flags       []Flag
}

func (c Command) Run(ctx *Context) error {
	c.Flags = append(
		c.Flags,
		helpFlag{"show help"},
	)
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
	if firstFlagIndex > -1 {
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
	context := NewContext(ctx.App, set, ctx.globalSet)
	if checkCommandHelp(context, c.Name) {
		return nil
	}
	c.Action(context)
	return nil
}

func (c Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}
