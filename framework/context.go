package framework

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Context struct {
	rsp http.ResponseWriter
	req *http.Request

	ctx context.Context

	hEntryIndex int
	hEntryList  []handlerFuncEntry
}

func NewContext(rsp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		rsp:         rsp,
		req:         req,
		ctx:         req.Context(),
		hEntryIndex: -1,
		hEntryList:  nil,
	}
}

var _ context.Context = &Context{}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// method to access internal variable.
// TODO: wrap those method better
func (c *Context) BaseContext() context.Context {
	return c.req.Context()
}

func (c *Context) Request() *http.Request {
	return c.req
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.rsp
}

func (c *Context) SetHandlerList(list []handlerFuncEntry) {
	c.hEntryList = list
}

func (c *Context) NextHandler() error {
	c.hEntryIndex++
	if c.hEntryIndex < len(c.hEntryList) {
		entry := c.hEntryList[c.hEntryIndex]

		fmt.Printf("--> [%d/%d] forward to [%s]\n",
			c.hEntryIndex+1, len(c.hEntryList), entry.funName)

		if err := entry.fun(c); err != nil {
			return err
		}
	}

	return nil
}
