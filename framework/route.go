package framework

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

func (core *Core) registerHelper(method, url string, fun HandlerFunc) {
	fName := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	fEntry := handlerFuncEntry{
		funName: fName + "()",
		pattern: url,
		fun:     fun,
	}

	if _, ok := core.routers[method]; !ok {
		// lazy allocate route trie tree
		core.routers[method] = newRouteTree(method)
	}

	if err := core.routers[method].addRoute(url, fEntry); err != nil {
		fmt.Println()
		log.Fatalf("add %s route Fatal with:[%s]\n", method, err.Error())
	}
}

func (core *Core) GetRegisterFunc(url string, fun HandlerFunc) {
	core.registerHelper("GET", url, fun)
}

func (core *Core) PostRegisterFunc(url string, fun HandlerFunc) {
	core.registerHelper("POST", url, fun)
}

func (core *Core) PutRegisterFunc(url string, fun HandlerFunc) {
	core.registerHelper("PUT", url, fun)
}

func (core *Core) DeleteRegisterFunc(url string, fun HandlerFunc) {
	core.registerHelper("DELETE", url, fun)
}

func (core *Core) getHandlerEntry(method, url string) *handlerFuncEntry {
	router, ok := core.routers[method]
	if !ok {
		return nil
	}

	return router.FindHandlerEntry(url)
}
