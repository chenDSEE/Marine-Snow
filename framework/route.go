package framework

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

func (core *Core) registerHelper(method, url string, fun []HandlerFunc) {
	fEntryList := make([]handlerFuncEntry, 0, len(fun)+len(core.middlewares))

	/* append HandlerFunc from core.middlewares */
	for _, f := range core.middlewares {
		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		fEntry := handlerFuncEntry{
			funName: fName + "()",
			pattern: url,
			fun:     f,
		}
		fEntryList = append(fEntryList, fEntry)
	}

	/* append HandlerFunc from parameters */
	for _, f := range fun {
		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		fEntry := handlerFuncEntry{
			funName: fName + "()",
			pattern: url,
			fun:     f,
		}
		fEntryList = append(fEntryList, fEntry)
	}

	if _, ok := core.routers[method]; !ok {
		// lazy allocate route trie tree
		core.routers[method] = newRouteTree(method)
	}

	if err := core.routers[method].addRoute(url, fEntryList); err != nil {
		fmt.Println()
		log.Fatalf("add %s route Fatal with:[%s]\n", method, err.Error())
	}
}

func (core *Core) GetRegisterFunc(url string, fun ...HandlerFunc) {
	core.registerHelper("GET", url, fun)
}

func (core *Core) PostRegisterFunc(url string, fun ...HandlerFunc) {
	core.registerHelper("POST", url, fun)
}

func (core *Core) PutRegisterFunc(url string, fun ...HandlerFunc) {
	core.registerHelper("PUT", url, fun)
}

func (core *Core) DeleteRegisterFunc(url string, fun ...HandlerFunc) {
	core.registerHelper("DELETE", url, fun)
}

func (core *Core) getHandlerEntryList(method, url string) ([]handlerFuncEntry, string) {
	router, ok := core.routers[method]
	if !ok {
		return nil, ""
	}

	return router.FindHandlerEntryList(url)
}
