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

type Args []string

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

func (c *Context) StringSlice(name string) []string {
	return c.lookupStringSlice(name, c.flagSet)
}

func (c *Context) IntSlice(name string) []int {
	return c.lookupIntSlice(name, c.flagSet)
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

func (c *Context) GlobalStringSlice(name string) []string {
	return c.lookupStringSlice(name, c.globalSet)
}

func (c *Context) GlobalIntSlice(name string) []int {
	return c.lookupIntSlice(name, c.globalSet)
}

func (c *Context) Args() Args {
	args := Args(c.flagSet.Args())
	return args
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

func (c *Context) lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*StringSlice)).Value()

	}
	return nil
}

func (c *Context) lookupIntSlice(name string, set *flag.FlagSet) []int {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*IntSlice)).Value()
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

func (a Args) Get(n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

func (a Args) First() string {
	return a.Get(0)
}

func (a Args) Tail() []string {
	if len(a) >= 2 {
		return []string(a)[1:]
	}
	return []string{}
}

func (a Args) Present() bool {
	return len(a) != 0
}
