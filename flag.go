package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

const defaultPlaceholder = "value"

var BashCompletionFlag Flag = BoolFlag{
	Name:   "generate-bash-completion",
	Hidden: true,
}

var VersionFlag Flag = BoolFlag{
	Name:  "version, v",
	Usage: "print the version",
}

var HelpFlag Flag = BoolFlag{
	Name:  "help, h",
	Usage: "show help",
}

var FlagStringer FlagStringFunc = stringifyFlag

var FlagNamePrefixer FlagNamePrefixFunc = prefixedNames

var FlagEnvHinter FlagEnvHintFunc = withEnvHint

var FlagFileHinter FlagFileHintFunc = withFileHint

type FlagsByName []Flag

func (f FlagsByName) Len() int {
	return len(f)
}

func (f FlagsByName) Less(i, j int) bool {
	return lexicographicLess(f[i].GetName(), f[j].GetName())
}

func (f FlagsByName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Flag interface {
	fmt.Stringer
	Apply(*flag.FlagSet)
	GetName() string
}

type RequiredFlag interface {
	Flag
	IsRequired() bool
}

type DocGenerationFlag interface {
	Flag
	TakesValue() bool
	GetUsage() string
	GetValue() string
}

type errorableFlag interface {
	Flag
	ApplyWithError(*flag.FlagSet) error
}

func flagSet(name string, flags []Flag) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	for _, f := range flags {
		if ef, ok := f.(errorableFlag); ok {
			if err := ef.ApplyWithError(set); err != nil {
				return nil, err
			}
		} else {
			f.Apply(set)
		}
	}
	set.SetOutput(ioutil.Discard)
	return set, nil
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

func visibleFlags(fl []Flag) []Flag {
	var visible []Flag
	for _, f := range fl {
		field := flagValue(f).FieldByName("Hidden")
		if !field.IsValid() || !field.Bool() {
			visible = append(visible, f)
		}
	}
	return visible
}

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func unquoteUsage(usage string) (string, string) {
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name := usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break
		}
	}
	return "", usage
}

func prefixedNames(fullName, placeholder string) string {
	var prefixed string
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += prefixFor(name) + name
		if placeholder != "" {
			prefixed += " " + placeholder
		}
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return prefixed
}

func withEnvHint(envVar, str string) string {
	envText := ""
	if envVar != "" {
		prefix := "$"
		suffix := ""
		sep := ", $"
		if runtime.GOOS == "windows" {
			prefix = "%"
			suffix = "%"
			sep = "%, %"
		}
		envText = " [" + prefix + strings.Join(strings.Split(envVar, ","), sep) + suffix + "]"
	}
	return str + envText
}

func withFileHint(filePath, str string) string {
	fileText := ""
	if filePath != "" {
		fileText = fmt.Sprintf(" [%s]", filePath)
	}
	return str + fileText
}

func flagValue(f Flag) reflect.Value {
	fv := reflect.ValueOf(f)
	for fv.Kind() == reflect.Ptr {
		fv = reflect.Indirect(fv)
	}
	return fv
}

func stringifyFlag(f Flag) string {
	fv := flagValue(f)
	switch f.(type) {
	case IntSliceFlag:
		return FlagFileHinter(
			fv.FieldByName("FilePath").String(),
			FlagEnvHinter(
				fv.FieldByName("EnvVar").String(),
				stringifyIntSliceFlag(f.(IntSliceFlag)),
			),
		)
	case Int64SliceFlag:
		return FlagFileHinter(
			fv.FieldByName("FilePath").String(),
			FlagEnvHinter(
				fv.FieldByName("EnvVar").String(),
				stringifyInt64SliceFlag(f.(Int64SliceFlag)),
			),
		)
	case StringSliceFlag:
		return FlagFileHinter(
			fv.FieldByName("FilePath").String(),
			FlagEnvHinter(
				fv.FieldByName("EnvVar").String(),
				stringifyStringSliceFlag(f.(StringSliceFlag)),
			),
		)
	}
	placeholder, usage := unquoteUsage(fv.FieldByName("Usage").String())
	needsPlaceholder := false
	defaultValueString := ""
	if val := fv.FieldByName("Value"); val.IsValid() {
		needsPlaceholder = true
		defaultValueString = fmt.Sprintf(" (default: %v)", val.Interface())

		if val.Kind() == reflect.String && val.String() != "" {
			defaultValueString = fmt.Sprintf(" (default: %q)", val.String())
		}
	}
	if defaultValueString == " (default: )" {
		defaultValueString = ""
	}
	if needsPlaceholder && placeholder == "" {
		placeholder = defaultPlaceholder
	}
	usageWithDefault := strings.TrimSpace(usage + defaultValueString)
	return FlagFileHinter(
		fv.FieldByName("FilePath").String(),
		FlagEnvHinter(
			fv.FieldByName("EnvVar").String(),
			FlagNamePrefixer(fv.FieldByName("Name").String(), placeholder)+"\t"+usageWithDefault,
		),
	)
}

func stringifyIntSliceFlag(f IntSliceFlag) string {
	var defaultVals []string
	if f.Value != nil && len(f.Value.Value()) > 0 {
		for _, i := range f.Value.Value() {
			defaultVals = append(defaultVals, strconv.Itoa(i))
		}
	}
	return stringifySliceFlag(f.Usage, f.Name, defaultVals)
}

func stringifyInt64SliceFlag(f Int64SliceFlag) string {
	var defaultVals []string
	if f.Value != nil && len(f.Value.Value()) > 0 {
		for _, i := range f.Value.Value() {
			defaultVals = append(defaultVals, strconv.FormatInt(i, 10))
		}
	}
	return stringifySliceFlag(f.Usage, f.Name, defaultVals)
}

func stringifyStringSliceFlag(f StringSliceFlag) string {
	var defaultVals []string
	if f.Value != nil && len(f.Value.Value()) > 0 {
		for _, s := range f.Value.Value() {
			if len(s) > 0 {
				defaultVals = append(defaultVals, strconv.Quote(s))
			}
		}
	}
	return stringifySliceFlag(f.Usage, f.Name, defaultVals)
}

func stringifySliceFlag(usage, name string, defaultVals []string) string {
	placeholder, usage := unquoteUsage(usage)
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}

	defaultVal := ""
	if len(defaultVals) > 0 {
		defaultVal = fmt.Sprintf(" (default: %s)", strings.Join(defaultVals, ", "))
	}

	usageWithDefault := strings.TrimSpace(usage + defaultVal)
	return FlagNamePrefixer(name, placeholder) + "\t" + usageWithDefault
}

func flagFromFileEnv(filePath, envName string) (val string, ok bool) {
	for _, envVar := range strings.Split(envName, ",") {
		envVar = strings.TrimSpace(envVar)
		if envVal, ok := syscall.Getenv(envVar); ok {
			return envVal, true
		}
	}
	for _, fileVar := range strings.Split(filePath, ",") {
		if fileVar != "" {
			if data, err := ioutil.ReadFile(fileVar); err == nil {
				return string(data), true
			}
		}
	}
	return "", false
}
