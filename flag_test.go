package cli

import (
	"testing"
)

func TestBoolFlagHelpOutput(t *testing.T) {
	var testCases = []struct {
		name     string
		expected string
	}{
		{"help", "--help\t"},
		{"h", "-h\t"},
	}

	for _, tc := range testCases {
		flag := BoolFlag{
			Name: tc.name,
		}
		output := flag.String()
		if output != tc.expected {
			t.Errorf("%s does not match %s", output, tc.expected)
		}
	}
}

func TestStringFlagHelpOutput(t *testing.T) {
	var testCases = []struct {
		name     string
		expected string
	}{
		{"help", "--help ''\t"},
		{"h", "-h ''\t"},
	}
	for _, tc := range testCases {
		flag := StringFlag{
			Name: tc.name,
		}
		output := flag.String()
		if output != tc.expected {
			t.Errorf("%s does not match %s", output, tc.expected)
		}
	}
}

func TestIntFlagHelpOutput(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{"help", "--help '0'\t"},
		{"h", "-h '0'\t"},
	}

	for _, test := range testCases {
		flag := IntFlag{
			Name: test.name,
		}
		output := flag.String()
		if output != test.expected {
			t.Errorf("%s does not match %s", output, test.expected)
		}
	}
}
