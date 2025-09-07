package cli

import (
	"fmt"
	"io/ioutil"
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
	err := set.Parse(ctx.Args()[1:])

	if err != nil {
		fmt.Println("Incorrect Usage.")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return err
	}
	if err := normalizeFlags(c.Flags, set); err != nil {
		fmt.Println(err)
		fmt.Println("")
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
