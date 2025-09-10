package cli

import "flag"

type StringFlag struct {
	Name        string
	Usage       string
	EnvVar      string
	FilePath    string
	Required    bool
	Hidden      bool
	TakesFile   bool
	Value       string
	Destination *string
}

func (f StringFlag) String() string {
	return FlagStringer(f)
}

func (f StringFlag) GetName() string {
	return f.Name
}

func (f StringFlag) IsRequired() bool {
	return f.Required
}

func (f StringFlag) TakesValue() bool {
	return true
}

func (f StringFlag) GetUsage() string {
	return f.Usage
}

func (f StringFlag) GetValue() string {
	return f.Value
}

func (f StringFlag) Apply(set *flag.FlagSet) {
	_ = f.ApplyWithError(set)
}

func (f StringFlag) ApplyWithError(set *flag.FlagSet) error {
	if envVal, ok := flagFromFileEnv(f.FilePath, f.EnvVar); ok {
		f.Value = envVal
	}
	eachName(f.Name, func(name string) {
		if f.Destination != nil {
			set.StringVar(f.Destination, name, f.Value, f.Usage)
			return
		}
		set.String(name, f.Value, f.Usage)
	})
	return nil
}

func (c *Context) String(name string) string {
	return lookupString(name, c.flagSet)
}

func (c *Context) GlobalString(name string) string {
	if fs := lookupGlobalFlagSet(name, c); fs != nil {
		return lookupString(name, fs)
	}
	return ""
}

func lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := f.Value.String(), error(nil)
		if err != nil {
			return ""
		}
		return parsed
	}
	return ""
}
