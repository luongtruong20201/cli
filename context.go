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

func (c *Context) Float64(name string) float64 {
	return lookupFloat64(name, c.flagSet)
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

type Args []string

func (c *Context) Args() Args {
	args := Args(c.flagSet.Args())
	return args
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

func lookupInt(name string, set *flag.FlagSet) int {
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

func lookupFloat64(name string, set *flag.FlagSet) float64 {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.ParseFloat(f.Value.String(), 64)
		if err != nil {
			return 0
		}
		return val
	}

	return 0
}

func lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	if f != nil {
		return f.Value.String()
	}

	return ""
}

func lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*StringSlice)).Value()

	}

	return nil
}

func lookupIntSlice(name string, set *flag.FlagSet) []int {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*IntSlice)).Value()

	}

	return nil
}

func lookupBool(name string, set *flag.FlagSet) bool {
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
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
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
