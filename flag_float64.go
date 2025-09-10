package cli

import (
	"flag"
	"fmt"
	"strconv"
)

type Float64Flag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       float64
	Destination *float64
}

func (f Float64Flag) String() string {
	return FlagStringer(f)
}

func (f Float64Flag) GetName() string {
	return f.Name
}

func (f Float64Flag) IsRequired() bool {
	return f.Required
}

func (f Float64Flag) TakesValue() bool {
	return true
}

func (f Float64Flag) GetUsage() string {
	return f.Usage
}

func (f Float64Flag) GetValue() string {
	return fmt.Sprintf("%f", f.Value)
}

func (c *Context) Float64(name string) float64 {
	return lookupFloat64(name, c.flagSet)
}

func (c *Context) GlobalFloat64(name string) float64 {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupFloat64(name, fs)
	}
	return 0
}

func (f Float64Flag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f Float64Flag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		envValFloat, err := strconv.ParseFloat(envVal, 10)
		if err != nil {
			return fmt.Errorf("could not parse %s as float64 value for flag %s: %s", envVal, f.Name, err)
		}

		f.Value = envValFloat
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.Float64Var(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.Float64(name, f.Value, f.Usage)
	})
	return nil
}

func lookupFloat64(name string, set *flag.FlagSet) float64 {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseFloat(f.Value.String(), 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}
