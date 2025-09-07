package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"
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
	return lookupInt(name, c.flagSet)
}

func (c *Context) Bool(name string) bool {
	return lookupBool(name, c.flagSet)
}

func (c *Context) String(name string) string {
	return lookupString(name, c.flagSet)
}

func (c *Context) StringSlice(name string) []string {
	return lookupStringSlice(name, c.flagSet)
}

func (c *Context) IntSlice(name string) []int {
	return lookupIntSlice(name, c.flagSet)
}

func (c *Context) GlobalInt(name string) int {
	return lookupInt(name, c.globalSet)
}

func (c *Context) GlobalBool(name string) bool {
	return lookupBool(name, c.globalSet)
}

func (c *Context) GlobalString(name string) string {
	return lookupString(name, c.globalSet)
}

func (c *Context) GlobalStringSlice(name string) []string {
	return lookupStringSlice(name, c.globalSet)
}

func (c *Context) GlobalIntSlice(name string) []int {
	return lookupIntSlice(name, c.globalSet)
}

func (c *Context) Args() []string {
	return c.flagSet.Args()
}

func lookupInt(name string, set *flag.FlagSet) int {
	if f := set.Lookup(name); f != nil {
		val, err := strconv.Atoi(f.Value.String())
		if err != nil {
			return 0
		}
		return val
	}
	return 0
}

func lookupString(name string, set *flag.FlagSet) string {
	if f := set.Lookup(name); f != nil {
		return f.Value.String()
	}
	return ""
}

func lookupStringSlice(name string, set *flag.FlagSet) []string {
	if f := set.Lookup(name); f != nil {
		return (f.Value.(*StringSlice)).Value()

	}
	return nil
}

func lookupIntSlice(name string, set *flag.FlagSet) []int {
	if f := set.Lookup(name); f != nil {
		return (f.Value.(*IntSlice)).Value()
	}
	return nil
}

func lookupBool(name string, set *flag.FlagSet) bool {
	if f := set.Lookup(name); f != nil {
		val, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return val
	}
	return false
}

func (c *Context) GetArg(n int) string {
	args := c.Args()
	if len(args) < n {
		return args[n]
	}
	return ""
}

func (c *Context) FirstArg() string {
	return c.GetArg(0)
}

func normalizeFlags(flags []Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.getName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name := strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			set.Set(name, ff.Value.String())
		}
	}
	return nil
}
