package framework

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	rsp http.ResponseWriter
	req *http.Request

	ctx context.Context

	hEntryIndex int
	hEntryList  []handlerFuncEntry

	queryAll   sync.Once
	queryTable map[string][]string
}

func NewContext(rsp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		rsp:         rsp,
		req:         req,
		ctx:         req.Context(),
		hEntryIndex: -1,
		hEntryList:  nil,
		queryTable:  map[string][]string{},
	}
}

var _ context.Context = &Context{}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

// method to access internal variable.
// TODO: wrap those method better
func (ctx *Context) BaseContext() context.Context {
	return ctx.req.Context()
}

func (ctx *Context) Request() *http.Request {
	return ctx.req
}

func (ctx *Context) ResponseWriter() http.ResponseWriter {
	return ctx.rsp
}

func (ctx *Context) SetHandlerList(list []handlerFuncEntry) {
	ctx.hEntryList = list
}

func (ctx *Context) NextHandler() error {
	ctx.hEntryIndex++
	if ctx.hEntryIndex < len(ctx.hEntryList) {
		entry := ctx.hEntryList[ctx.hEntryIndex]

		fmt.Printf("--> [%d/%d] forward to [%s]\n",
			ctx.hEntryIndex+1, len(ctx.hEntryList), entry.funName)

		if err := entry.fun(ctx); err != nil {
			return err
		}
	}

	return nil
}

var _ IRequest = &Context{}

func (ctx *Context) QueryAll() map[string][]string {
	ctx.queryAll.Do(func() {
		if ctx.req != nil {
			// every call for ctx.req.URL.Query() will parse URL and create a new map
			ctx.queryTable = ctx.req.URL.Query()
		}
	})

	return ctx.queryTable
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.Atoi(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseInt(valList[0], 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseBool(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList[0], true
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList, true
}

func (ctx *Context) Query(key string) interface{} {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return nil
	}

	return valList[0]
}
