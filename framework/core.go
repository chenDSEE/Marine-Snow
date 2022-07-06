package framework

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(c *Context) error

type Core struct {
	name   string
	router map[string]HandlerFunc
}

func NewCore() *Core {
	return &Core{
		name:   "MarineSnow Core demo",
		router: make(map[string]HandlerFunc),
	}
}

// as a HTTP URL router
func (core *Core) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.String())
	fun := core.router[req.URL.String()]
	if fun == nil {
		fmt.Printf("can not find any handler for [%s]\n", req.URL.String())
		return
	}

	ctx := NewContext(rsp, req)

	_ = fun(ctx) // TODO: print the accessed URL and handler name before jump to fun()
}

func (c Core) RegisterHandlerFunc(url string, fun HandlerFunc) {
	c.router[url] = fun
}
