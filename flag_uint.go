package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type UintFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       uint
	Destination *uint
}

func (f UintFlag) String() string {
	return FlagStringer(f)
}

func (f UintFlag) GetName() string {
	return f.Name
}

func (f UintFlag) IsRequired() bool {
	return f.Required
}

func (f UintFlag) TakesValue() bool {
	return true
}

func (f UintFlag) GetUsage() string {
	return f.Usage
}

func (f UintFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f UintFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		envValInt, err := strconv.ParseUint(envVal, 0, 64)
		if err != nil {
			return fmt.Errorf("could not parse %s as uint value for flag %s: %s", envVal, f.Name, err)
		}
		f.Value = uint(envValInt)
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.UintVar(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.Uint(name, f.Value, f.Usage)
	})
	return nil
}

func (f UintFlag) GetValue() string {
	return fmt.Sprintf("%d", f.Value)
}

func (c *Context) Uint(name string) uint {
	return lookupUint(name, c.flagSet)
}

func (c *Context) GlobalUint(name string) uint {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupUint(name, fs)
	}
	return 0
}

func lookupUint(name string, set *flag.FlagSet) uint {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseUint(f.Value.String(), 0, 64)
		if err != nil {
			return 0
		}
		return uint(parsed)
	}
	return 0
}
