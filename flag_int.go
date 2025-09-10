package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type IntFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       int
	Destination *int
}

func (f IntFlag) String() string {
	return FlagStringer(f)
}

func (f IntFlag) GetName() string {
	return f.Name
}

func (f IntFlag) IsRequired() bool {
	return f.Required
}

func (f IntFlag) TakesValue() bool {
	return true
}

func (f IntFlag) GetUsage() string {
	return f.Usage
}

func (f IntFlag) GetValue() string {
	return fmt.Sprintf("%d", f.Value)
}

func (f IntFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f IntFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		envValInt, err := strconv.ParseInt(envVal, 0, 64)
		if err != nil {
			return fmt.Errorf("could not parse %s as int value for flag %s: %s", envVal, f.Name, err)
		}
		f.Value = int(envValInt)
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.IntVar(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.Int(name, f.Value, f.Usage)
	})
	return nil
}

func (c *Context) Int(name string) int {
	return lookupInt(name, c.flagSet)
}

func (c *Context) GlobalInt(name string) int {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupInt(name, fs)
	}
	return 0
}

func lookupInt(name string, set *flag.FlagSet) int {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseInt(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		return int(parsed)
	}
	return 0
}
