package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"
)

var AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [global options] command [command options] [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

var CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   command {{.Name}} [command options] [arguments...]

DESCRIPTION:
   {{.Description}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

var ManPageTemplate = `.\"             -*-Nroff-*-
.\"
.TH "{{.Name}}" 1 "{{.Compiled.Day}} {{.Compiled.Month}} {{.Compiled.Year}}" "" ""
.SH NAME
{{.Name}} \- {{.Usage}}
.SH SYNOPSIS
.B {{.Name}}
.nf
command {{.Name}} [command options] [arguments...]
.fi
.SH COMMANDS
.nf
{{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
{{end}}
.fi
.SH OPTIONS
.nf
{{range .Flags}}{{.}}
{{end}}
.fi
.SH VERSION
.B {{.Version}}
.SH AUTHOR
.B {{.Name}} was written by {{.Author}} <{{.Email}}>
`

var helpCommand = Command{
	Name:      "help",
	ShortName: "h",
	Usage:     "Shows a list of commands or help for one command",
	Action: func(c *Context) {
		args := c.Args()
		if len(args) > 0 {
			ShowCommandHelp(c, args[0])
		} else {
			ShowAppHelp(c)
		}
	},
}

func ShowAppHelp(c *Context) {
	printHelp(AppHelpTemplate, c.App)
}

func GenerateManPage(c *Context) {
	printHelp(ManPageTemplate, c.App)
}

func ShowCommandHelp(c *Context, command string) {
	for _, c := range c.App.Commands {
		if c.HasName(command) {
			printHelp(CommandHelpTemplate, c)
			return
		}
	}

	fmt.Printf("No help topic for '%v'\n", command)
}

func ShowVersion(c *Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
}

func printHelp(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
	w.Flush()
}

func checkVersion(c *Context) bool {
	if c.GlobalBool("version") {
		ShowVersion(c)
		return true
	}
	return false
}

func checkHelp(c *Context) bool {
	if c.GlobalBool("h") || c.GlobalBool("help") {
		ShowAppHelp(c)
		return true
	}
	return false
}

func checkCommandHelp(c *Context, name string) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowCommandHelp(c, name)
		return true
	}
	return false
}
