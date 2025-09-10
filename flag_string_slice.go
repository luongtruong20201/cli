package cli

import (
	"flag"
	"fmt"
	"strings"
)

type StringSlice []string

func (f *StringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *StringSlice) String() string {
	return strings.Join(*f, ",")
}

func (f *StringSlice) Value() []string {
	return *f
}

func (f *StringSlice) Get() interface{} {
	return *f
}

type StringSliceFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Required  bool
	Hidden    bool
	TakesFile bool
	Value     *StringSlice
}

func (f StringSliceFlag) String() string {
	return FlagStringer(f)
}

func (f StringSliceFlag) GetName() string {
	return f.Name
}

func (f StringSliceFlag) IsRequired() bool {
	return f.Required
}

func (f StringSliceFlag) TakesValue() bool {
	return true
}

func (f StringSliceFlag) GetUsage() string {
	return f.Usage
}

func (f StringSliceFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

func (f StringSliceFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f StringSliceFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		newVal := &StringSlice{}
		for _, s := range strings.Split(envVal, ",") {
			s = strings.TrimSpace(s)
			if err := newVal.Set(s); err != nil {
				return fmt.Errorf("could not parse %s as string value for flag %s: %s", envVal, f.Name, err)
			}
		}
		if f.Value == nil {
			f.Value = newVal
		} else {
			*f.Value = *newVal
		}
	}
	eachName(f.Name, func(name string) {
		if f.Value == nil {
			f.Value = &StringSlice{}
		}
		set.Var(f.Value, name, f.Usage)
	})
	return nil
}

func (c *Context) StringSlice(name string) []string {
	return lookupStringSlice(name, c.flagSet)
}

func (c *Context) GlobalStringSlice(name string) []string {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupStringSlice(name, fs)
	}
	return nil
}

func lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	if f != nil {
		value, ok := f.Value.(*StringSlice)
		if !ok {
			return nil
		}
		slice := value.Value()
		var defaultVal []string
		for _, v := range strings.Split(f.DefValue, ",") {
			defaultVal = append(defaultVal, v)
		}
		if !isStringSliceEqual(slice, defaultVal) {
			for _, v := range defaultVal {
				slice = removeFromStringSlice(slice, v)
			}
		}
		return slice
	}
	return nil
}

func removeFromStringSlice(slice []string, val string) []string {
	for i, v := range slice {
		if v == val {
			ret := append([]string{}, slice[:i]...)
			ret = append(ret, slice[i+1:]...)
			return ret
		}
	}
	return slice
}

func isStringSliceEqual(newValue, defaultValue []string) bool {
	if (newValue == nil) != (defaultValue == nil) {
		return false
	}
	if len(newValue) != len(defaultValue) {
		return false
	}
	for i, v := range newValue {
		if v != defaultValue[i] {
			return false
		}
	}
	return true
}
