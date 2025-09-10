package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type Int64Flag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       int64
	Destination *int64
}

func (f Int64Flag) String() string {
	return FlagStringer(f)
}

func (f Int64Flag) GetName() string {
	return f.Name
}

func (f Int64Flag) IsRequired() bool {
	return f.Required
}

func (f Int64Flag) TakesValue() bool {
	return true
}

func (f Int64Flag) GetUsage() string {
	return f.Usage
}

func (f Int64Flag) GetValue() string {
	return fmt.Sprintf("%d", f.Value)
}

func (f Int64Flag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f Int64Flag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		envValInt, err := strconv.ParseInt(envVal, 0, 64)
		if err != nil {
			return fmt.Errorf("could not parse %s as int value for flag %s: %s", envVal, f.Name, err)
		}

		f.Value = envValInt
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.Int64Var(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.Int64(name, f.Value, f.Usage)
	})
	return nil
}

func (c *Context) Int64(name string) int64 {
	return lookupInt64(name, c.flagSet)
}

func (c *Context) GlobalInt64(name string) int64 {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupInt64(name, fs)
	}
	return 0
}

func lookupInt64(name string, set *flag.FlagSet) int64 {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseInt(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}
