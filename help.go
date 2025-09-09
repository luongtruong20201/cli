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
   {{.Name}} {{ if .Flags }}[global options] {{ end }}command{{ if .Flags }} [command options]{{ end }} [arguments...]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{ if .Flags }}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

var CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   command {{.Name}}{{ if .Flags }} [command options]{{ end }} [arguments...]

DESCRIPTION:
   {{.Description}}{{ if .Flags }}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

var SubcommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} command{{ if .Flags }} [command options]{{ end }} [arguments...]

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{ if .Flags }}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

var helpCommand = Command{
	Name:      "help",
	ShortName: "h",
	Usage:     "Shows a list of commands or help for one command",
	Action: func(c *Context) {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowAppHelp(c)
		}
	},
}

var helpSubcommand = Command{
	Name:      "help",
	ShortName: "h",
	Usage:     "Shows a list of commands or help for one command",
	Action: func(c *Context) {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowSubcommandHelp(c)
		}
	},
}

var HelpPrinter = printHelp

func ShowAppHelp(c *Context) {
	HelpPrinter(AppHelpTemplate, c.App)
}

func DefaultAppComplete(c *Context) {
	for _, command := range c.App.Commands {
		fmt.Println(command.Name)
		if command.ShortName != "" {
			fmt.Println(command.ShortName)
		}
	}
}

func ShowCommandHelp(c *Context, command string) {
	for _, c := range c.App.Commands {
		if c.HasName(command) {
			HelpPrinter(CommandHelpTemplate, c)
			return
		}
	}
	if c.App.CommandNotFound != nil {
		c.App.CommandNotFound(c, command)
	} else {
		fmt.Printf("No help topic for '%v'\n", command)
	}
}

func ShowSubcommandHelp(c *Context) {
	HelpPrinter(SubcommandHelpTemplate, c.App)
}

func ShowVersion(c *Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
}

func ShowCompletions(c *Context) {
	a := c.App
	if a != nil && a.BashComplete != nil {
		a.BashComplete(c)
	}
}

func ShowCommandCompletions(ctx *Context, command string) {
	c := ctx.App.Command(command)
	if c != nil && c.BashComplete != nil {
		c.BashComplete(ctx)
	}
}

func printHelp(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
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

func checkSubcommandHelp(c *Context) bool {
	if c.GlobalBool("h") || c.GlobalBool("help") {
		ShowSubcommandHelp(c)
		return true
	}

	return false
}

func checkCompletions(c *Context) bool {
	if c.GlobalBool(BashCompletionFlag.Name) && c.App.EnableBashCompletion {
		ShowCompletions(c)
		return true
	}
	return false
}

func checkCommandCompletions(c *Context, name string) bool {
	if c.Bool(BashCompletionFlag.Name) && c.App.EnableBashCompletion {
		ShowCommandCompletions(c, name)
		return true
	}
	return false
}
