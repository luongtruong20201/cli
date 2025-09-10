package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type Uint64Flag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       uint64
	Destination *uint64
}

func (f Uint64Flag) String() string {
	return FlagStringer(f)
}

func (f Uint64Flag) GetName() string {
	return f.Name
}

func (f Uint64Flag) IsRequired() bool {
	return f.Required
}

func (f Uint64Flag) TakesValue() bool {
	return true
}

func (f Uint64Flag) GetUsage() string {
	return f.Usage
}

func (f Uint64Flag) GetValue() string {
	return fmt.Sprintf("%d", f.Value)
}

func (f Uint64Flag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f Uint64Flag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		envValInt, err := strconv.ParseUint(envVal, 0, 64)
		if err != nil {
			return fmt.Errorf("could not parse %s as uint64 value for flag %s: %s", envVal, f.Name, err)
		}
		f.Value = envValInt
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.Uint64Var(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.Uint64(name, f.Value, f.Usage)
	})
	return nil
}

func (c *Context) Uint64(name string) uint64 {
	return lookupUint64(name, c.flagSet)
}

func (c *Context) GlobalUint64(name string) uint64 {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupUint64(name, fs)
	}
	return 0
}

func lookupUint64(name string, set *flag.FlagSet) uint64 {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseUint(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}
