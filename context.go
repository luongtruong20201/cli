package cli

import "flag"

type Context struct {
	App           *App
	Command       Command
	shellComplete bool
	flagSet       *flag.FlagSet
	setFlags      map[string]bool
	parentContext *Context
}

func NewContext(app *App, set *flag.FlagSet, parentCtx *Context) *Context {
	c := &Context{App: app, flagSet: set, parentContext: parentCtx}
	if parentCtx != nil {
		c.shellComplete = parentCtx.shellComplete
	}
	return c
}

func (c *Context) NumFlags() int {
	return c.flagSet.NFlag()
}

func (c *Context) Set(name, value string) error {
	c.setFlags = nil
	return c.flagSet.Set(name, value)
}

func (c *Context) GlobalSet(name, value string) error {
	globalContext(c).setFlags = nil
	return globalContext(c).flagSet.Set(name, value)
}

func globalContext(ctx *Context) *Context {
	if ctx == nil {
		return nil
	}
	for {
		if ctx.parentContext == nil {
			return ctx
		}
		ctx = ctx.parentContext
	}
}

func lookupGlobalFlagSet(name string, ctx *Context) *flag.FlagSet {
	if ctx.parentContext != nil {
		ctx = ctx.parentContext
	}
	for ; ctx != nil; ctx = ctx.parentContext {
		if f := ctx.flagSet.Lookup(name); f != nil {
			return ctx.flagSet
		}
	}
	return nil
}
