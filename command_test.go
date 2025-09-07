package cli_test

import (
	"cli"
	"flag"
	"testing"
)

func TestCommandDoNotIgoreFlags(t *testing.T) {
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	tests := []string{"blah", "blah", "-break"}
	set.Parse(tests)
	c := cli.NewContext(app, set, set)

	command := cli.Command{
		Name:        "test-cmd",
		ShortName:   "tc",
		Usage:       "this is for testing",
		Description: "testing",
		Action:      func(context *cli.Context) {},
	}
	err := command.Run(c)
	expect(t, err.Error(), "flag provided but not defined: -break")
}
