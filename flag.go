package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type Flag interface {
	fmt.Stringer
	Apply(*flag.FlagSet)
	GetName() string
}

func flagSet(name string, flags []Flag) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		f.Apply(set)
	}
	return set
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

type StringSlice []string

func (f *StringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *StringSlice) String() string {
	return fmt.Sprintf("%s", *f)
}

func (f *StringSlice) Value() []string {
	return *f
}

type StringSliceFlag struct {
	Name  string
	Value *StringSlice
	Usage string
}

func (f StringSliceFlag) String() string {
	return fmt.Sprintf("%s '%v'\t%v", prefixFor(f.Name), f.Name, "-"+f.Name+" option -"+f.Name+" option")
}

func (f StringSliceFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Usage)
	})
}

func (f StringSliceFlag) GetName() string {
	return f.Name
}

type IntSlice []int

func (f *IntSlice) Set(value string) error {

	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	} else {
		*f = append(*f, tmp)
	}
	return nil
}

func (f *IntSlice) String() string {
	return fmt.Sprintf("%d", *f)
}

func (f *IntSlice) Value() []int {
	return *f
}

type IntSliceFlag struct {
	Name  string
	Value *IntSlice
	Usage string
}

func (f IntSliceFlag) String() string {
	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
	pref := prefixFor(firstName)
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage)
}

func (f IntSliceFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Usage)
	})
}

func (f IntSliceFlag) GetName() string {
	return f.Name
}

type BoolFlag struct {
	Name  string
	Usage string
}

func (f BoolFlag) String() string {
	return fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Usage)
}

func (f BoolFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Bool(name, false, f.Usage)
	})
}

func (f BoolFlag) GetName() string {
	return f.Name
}

type StringFlag struct {
	Name  string
	Value string
	Usage string
}

func (f StringFlag) String() string {
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

func (f StringFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.String(name, f.Value, f.Usage)
	})
}

func (f StringFlag) GetName() string {
	return f.Name
}

type IntFlag struct {
	Name  string
	Value int
	Usage string
}

func (f IntFlag) String() string {
	return fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage)
}

func (f IntFlag) Apply(set *flag.FlagSet) {
	eachName(f.Name, func(name string) {
		set.Int(name, f.Value, f.Usage)
	})
}

func (f IntFlag) GetName() string {
	return f.Name
}

type helpFlag struct {
	Usage string
}

func (f helpFlag) String() string {
	return fmt.Sprintf("%s\t%v", prefixedNames("help, h"), f.Usage)
}

func (f helpFlag) Apply(set *flag.FlagSet) {
	eachName("help, h", func(name string) {
		set.Bool(name, false, f.Usage)
	})
}

func (f helpFlag) GetName() string {
	return "help"
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func prefixedNames(fullName string) (prefixed string) {
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += prefixFor(name) + name
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return
}
