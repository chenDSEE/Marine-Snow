package framework

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

type HandlerFunc func(c *Context) error

type Core struct {
	name   string
	router map[string]handlerFuncEntry
}

type handlerFuncEntry struct {
	funName string
	fun     HandlerFunc
}

func NewCore() *Core {
	return &Core{
		name:   "MarineSnow Core demo",
		router: make(map[string]handlerFuncEntry),
	}
}

// as a HTTP URL router
func (core *Core) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	funEntry, ok := core.router[url]
	if ok == false {
		fmt.Printf("can not find any handler for [%s]\n", req.URL.String())
		return
	}

	fmt.Printf("==> request[%s], forwarding to [%s]\n", url, funEntry.funName)
	ctx := NewContext(rsp, req)

	_ = funEntry.fun(ctx)
}

func (c Core) RegisterHandlerFunc(url string, fun HandlerFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	c.router[url] = handlerFuncEntry{
		funName: name + "()",
		fun:     fun,
	}
}
