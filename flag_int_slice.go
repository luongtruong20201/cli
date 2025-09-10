package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

func (f *IntSlice) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, tmp)
	return nil
}

func (f *IntSlice) String() string {
	slice := make([]string, len(*f))
	for i, v := range *f {
		slice[i] = strconv.Itoa(v)
	}
	return strings.Join(slice, ",")
}

func (f *IntSlice) Value() []int {
	return *f
}

func (f *IntSlice) Get() interface{} {
	return *f
}

type IntSliceFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	FilePath string
	Required bool
	Hidden   bool
	Value    *IntSlice
}

func (f IntSliceFlag) String() string {
	return FlagStringer(f)
}

func (f IntSliceFlag) GetName() string {
	return f.Name
}

func (f IntSliceFlag) IsRequired() bool {
	return f.Required
}

func (f IntSliceFlag) TakesValue() bool {
	return true
}

func (f IntSliceFlag) GetUsage() string {
	return f.Usage
}

func (f IntSliceFlag) GetValue() string {
	if f.Value != nil {
		return f.Value.String()
	}
	return ""
}

func (f IntSliceFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f IntSliceFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		newVal := &IntSlice{}
		for _, s := range strings.Split(envVal, ",") {
			s = strings.TrimSpace(s)
			if err := newVal.Set(s); err != nil {
				return fmt.Errorf("could not parse %s as int slice value for flag %s: %s", envVal, f.Name, err)
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
			f.Value = &IntSlice{}
		}
		set.Var(f.Value, name, f.Usage)
	})
	return nil
}

func (c *Context) IntSlice(name string) []int {
	return lookupIntSlice(name, c.flagSet)
}

func lookupIntSlice(name string, set *flag.FlagSet) []int {
	if f := set.Lookup(name); f != nil {
		value, ok := f.Value.(*IntSlice)
		if !ok {
			return nil
		}
		slice := value.Value()
		var defaultVal []int
		for _, v := range strings.Split(f.DefValue, ",") {
			if v != "" {
				intValue, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				defaultVal = append(defaultVal, intValue)
			}
		}
		if !isIntSliceEqual(slice, defaultVal) {
			for _, v := range defaultVal {
				slice = removeFromIntSlice(slice, v)
			}
		}
		return slice
	}
	return nil
}

func removeFromIntSlice(slice []int, val int) []int {
	for i, v := range slice {
		if v == val {
			ret := append([]int{}, slice[:i]...)
			ret = append(ret, slice[i+1:]...)
			return ret
		}
	}
	return slice
}

func isIntSliceEqual(newValue, defaultValue []int) bool {
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
