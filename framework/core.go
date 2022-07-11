package framework

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context) error

type Core struct {
	name        string
	middlewares []HandlerFunc
	routers     map[string]*rTree // router tire tree
}

type handlerFuncEntry struct {
	funName string
	pattern string
	fun     HandlerFunc
}

func NewCore() *Core {
	core := &Core{
		name:        "MarineSnow Core demo",
		middlewares: make([]HandlerFunc, 0),
		routers:     make(map[string]*rTree),
	}

	return core
}

// as a HTTP URL router
func (core *Core) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	method := strings.ToUpper(req.Method)
	fEntryList, pattern := core.getHandlerEntryList(method, url)
	if fEntryList == nil || len(fEntryList) == 0 {
		fmt.Printf("can not find any handler for [%s:%s]\n", method, url)
		return
	}

	fmt.Printf("==> request[%s:%s], match [%s], start pipeline handle:\n",
		method, url, pattern)

	ctx := NewContext(rsp, req)
	ctx.SetHandlerList(fEntryList)

	if err := ctx.NextHandler(); err != nil {
		fmt.Printf("Catch an error: %s\n", err.Error())
		return
	}
}

func (core *Core) NewRouteGroup(prefix string) RouteGroup {
	return newPrefixGroup(core, prefix)
}

func (core *Core) DumpRoutes() {
	for _, router := range core.routers {
		router.printRouteTree()
		fmt.Println("")
	}
}

func (core *Core) DumpMethodRoute(method string) {
	method = strings.ToUpper(method)
	router, ok := core.routers[method]
	if !ok {
		return
	}

	router.printRouteTree()
}

func (core *Core) AppendDefaultMiddleware(middlewareFun ...HandlerFunc) {
	core.middlewares = append(core.middlewares, middlewareFun...)
}
