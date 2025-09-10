package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type BoolTFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Destination *bool
}

func (f BoolTFlag) String() string {
	return FlagStringer(f)
}

func (f BoolTFlag) GetName() string {
	return f.Name
}

func (f BoolTFlag) IsRequired() bool {
	return f.Required
}

func (f BoolTFlag) TakesValue() bool {
	return false
}

func (f BoolTFlag) GetUsage() string {
	return f.Usage
}

func (f BoolTFlag) GetValue() string {
	return ""
}

func (c *Context) BoolT(name string) bool {
	return lookupBoolT(name, c.flagSet)
}

func (c *Context) GlobalBoolT(name string) bool {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupBoolT(name, fs)
	}
	return false
}

func (f BoolTFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f BoolTFlag) ApplyWithError(set *flag.FlagSet) error {
	val := true
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

func lookupBoolT(name string, set *flag.FlagSet) bool {
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
