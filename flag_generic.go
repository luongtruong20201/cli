package cli

import (
	"flag"
	"fmt"
)

type Generic interface {
	Set(value string) error
	String() string
}

type GenericFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Required  bool
	Hidden    bool
	TakesFile bool
	Value     Generic
}

func (f GenericFlag) String() string {
	return FlagStringer(f)
}

func (f GenericFlag) GetName() string {
	return f.Name
}

func (f GenericFlag) IsRequired() bool {
	return f.Required
}

func (f GenericFlag) TakesValue() bool {
	return true
}

func (f GenericFlag) GetUsage() string {
	return f.Usage
}

func (f GenericFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

func (f GenericFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f GenericFlag) ApplyWithError(set *flag.FlagSet) error {
	val := f.Value
	if fileEnvVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		if err := val.Set(fileEnvVal); err != nil {
			return fmt.Errorf("could not parse %s as value for flag %s: %s", fileEnvVal, f.Name, err)
		}
	}
	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Usage)
	})
	return nil
}

func (c *Context) Generic(name string) interface{} {
	return lookupGeneric(name, c.flagSet)
}

func (c *Context) GlobalGeneric(name string) interface{} {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupGeneric(name, fs)
	}
	return nil
}

func lookupGeneric(name string, set *flag.FlagSet) interface{} {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := f.Value, error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
