package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Int64Slice []int64

func (f *Int64Slice) Set(value string) error {
	tmp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*f = append(*f, tmp)
	return nil
}

func (f *Int64Slice) String() string {
	slice := make([]string, len(*f))
	for i, v := range *f {
		slice[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(slice, ",")
}

func (f *Int64Slice) Value() []int64 {
	return *f
}

func (f *Int64Slice) Get() interface{} {
	return *f
}

type Int64SliceFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	FilePath string
	Required bool
	Hidden   bool
	Value    *Int64Slice
}

func (f Int64SliceFlag) String() string {
	return FlagStringer(f)
}

func (f Int64SliceFlag) GetName() string {
	return f.Name
}

func (f Int64SliceFlag) IsRequired() bool {
	return f.Required
}

func (f Int64SliceFlag) TakesValue() bool {
	return true
}

func (f Int64SliceFlag) GetUsage() string {
	return f.Usage
}

func (f Int64SliceFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

func (f Int64SliceFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f Int64SliceFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		newVal := &Int64Slice{}
		for _, s := range strings.Split(envVal, ",") {
			s = strings.TrimSpace(s)
			if err := newVal.Set(s); err != nil {
				return fmt.Errorf("could not parse %s as int64 slice value for flag %s: %s", envVal, f.Name, err)
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
			f.Value = &Int64Slice{}
		}
		set.Var(f.Value, name, f.Usage)
	})
	return nil
}

func (c *Context) Int64Slice(name string) []int64 {
	return lookupInt64Slice(name, c.flagSet)
}

func (c *Context) GlobalInt64Slice(name string) []int64 {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupInt64Slice(name, fs)
	}
	return nil
}

func lookupInt64Slice(name string, set *flag.FlagSet) []int64 {
	f := set.Lookup(name)
	if f != nil {
		value, ok := f.Value.(*Int64Slice)
		if !ok {
			return nil
		}
		parsed := value.Value()
		var defaultVal []int64
		for _, v := range strings.Split(f.DefValue, ",") {
			if v != "" {
				int64Value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(err)
				}
				defaultVal = append(defaultVal, int64Value)
			}
		}
		if !isInt64SliceEqual(parsed, defaultVal) {
			for _, v := range defaultVal {
				parsed = removeFromInt64Slice(parsed, v)
			}
		}
		return parsed
	}
	return nil
}

func removeFromInt64Slice(slice []int64, val int64) []int64 {
	for i, v := range slice {
		if v == val {
			ret := append([]int64{}, slice[:i]...)
			ret = append(ret, slice[i+1:]...)
			return ret
		}
	}
	return slice
}

func isInt64SliceEqual(newValue, defaultValue []int64) bool {
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
