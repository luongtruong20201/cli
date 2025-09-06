package cli

// var (
// 	Name     = os.Args[0]
// 	Usage    = "<No Description>"
// 	Version  = "0.0.0"
// 	Commands []Command
// 	Flags    []Flag
// 	Action   = ShowHelp
// )

// func Run(arguments []string) {
// 	set := flagSet(Flags)
// 	set.Parse(arguments[1:])
// 	context := NewContext(set, set)
// 	args := context.Args()
// 	if len(args) > 0 {
// 		name := args[0]
// 		for _, c := range append(Commands, HelpCommand) {
// 			if c.HasName(name) {
// 				locals := flagSet(c.Flags)
// 				locals.Parse(args[1:])
// 				c.Run(context)
// 				return
// 			}
// 		}
// 	}
// 	Action(context)
// }

// func CommandsWithDefaults() []Command {
// 	return append(append([]Command(nil), HelpCommand), Commands...)
// }

type Handler func(context *Context)
