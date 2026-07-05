package types

import (
	"context"
)

type Context struct {
	ctx context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}

func (c *Context) GetContext() context.Context {
	return c.ctx
}
