package cli

type BashCompleteFunc func(*Context)

type BeforeFunc func(*Context) error

type AfterFunc func(ctx *Context) error

type ActionFunc func(*Context) error

type CommandNotFoundFunc func(*Context, string)

type OnUsageErrorFunc func(*Context, error, bool) error

type ExitErrHandleFunc func(*Context, error)

type FlagStringFunc func(Flag) string

type FlagNamePrefixFunc func(string, string) string

type FlagEnvHintFunc func(string, string) string

type FlagFileHintFunc func(string, string) string
