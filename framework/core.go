package framework

import (
	"context"
	"fmt"
	"net/http"
)

type HandlerFunc func(c Context) error

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

	ctx := Context{
		Rsp: rsp,
		Req: req,
		Ctx: context.Background(),
	}

	fun(ctx)
}

func (c Core) RegisterHandlerFunc(url string, fun HandlerFunc) {
	c.router[url] = fun
}
