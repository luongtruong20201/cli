package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type BoolFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Destination *bool
}

func (f BoolFlag) String() string {
	return FlagStringer(f)
}

func (f BoolFlag) GetName() string {
	return f.Name
}

func (f BoolFlag) IsRequired() bool {
	return f.Required
}

func (f BoolFlag) TakesValue() bool {
	return false
}

func (f BoolFlag) GetUsage() string {
	return f.Usage
}

func (f BoolFlag) GetValue() string {
	return ""
}

func (c *Context) Bool(name string) bool {
	return lookupBool(name, c.flagSet)
}

func (c *Context) GlobalBool(name string) bool {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupBool(name, fs)
	}
	return false
}

func (f BoolFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f BoolFlag) ApplyWithError(set *flag.FlagSet) error {
	val := false
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		if envVal == "" {
			val = false
		} else {
			envValBool, err := strconv.ParseBool(envVal)
			if err != nil {
				return fmt.Errorf("could not parse %s as bool value for flag %s: %s", envVal, f.Name, err)
			}
			val = envValBool
		}
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.BoolVar(f.Destination, name, val, f.Usage)
			return
		}
		set.Bool(name, val, f.Usage)
	})
	return nil
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}
