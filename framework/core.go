package framework

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context) error

type Core struct {
	name    string
	routers map[string]*rTree // router tire tree
}

type handlerFuncEntry struct {
	funName string
	pattern string
	fun     HandlerFunc
}

func NewCore() *Core {
	core := &Core{
		name:    "MarineSnow Core demo",
		routers: make(map[string]*rTree),
	}

	return core
}

// as a HTTP URL router
func (core *Core) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	method := strings.ToUpper(req.Method)
	fEntry := core.getHandlerEntry(method, url)
	if fEntry == nil {
		fmt.Printf("can not find any handler for [%s:%s]\n", method, url)
		return
	}

	fmt.Printf("==> request[%s:%s], match [%s], forwarding to [%s]\n",
		method, url, fEntry.pattern, fEntry.funName)

	ctx := NewContext(rsp, req)

	_ = fEntry.fun(ctx)
}

func (core *Core) NewRouteGroup(prefix string) RouteGroup {
	return newPrefixGroup(core, prefix)
}
