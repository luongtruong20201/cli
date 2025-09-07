package cli

import (
	"flag"
	"strconv"
)

type Context struct {
	App       *App
	flagSet   *flag.FlagSet
	globalSet *flag.FlagSet
}

func NewContext(app *App, set *flag.FlagSet, globalSet *flag.FlagSet) *Context {
	return &Context{app, set, globalSet}
}

func (c *Context) Int(name string) int {
	return c.lookupInt(name, c.flagSet)
}

func (c *Context) Bool(name string) bool {
	return c.lookupBool(name, c.flagSet)
}

func (c *Context) String(name string) string {
	return c.lookupString(name, c.flagSet)
}

func (c *Context) StringSlice(name string) flag.Value {
	return c.lookupStringSlice(name, c.flagSet)
}

func (c *Context) GlobalInt(name string) int {
	return c.lookupInt(name, c.globalSet)
}

func (c *Context) GlobalBool(name string) bool {
	return c.lookupBool(name, c.globalSet)
}

func (c *Context) GlobalString(name string) string {
	return c.lookupString(name, c.globalSet)
}

func (c *Context) Args() []string {
	return c.flagSet.Args()
}

func (c *Context) lookupInt(name string, set *flag.FlagSet) int {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.Atoi(f.Value.String())
		if err != nil {
			return 0
		}
		return val
	}

	return 0
}

func (c *Context) lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	if f != nil {
		return f.Value.String()
	}

	return ""
}

func (c *Context) lookupStringSlice(name string, set *flag.FlagSet) flag.Value {
	if f := set.Lookup(name); f != nil {
		return f.Value
	}
	return nil
}

func (c *Context) lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return val
	}

	return false
}
