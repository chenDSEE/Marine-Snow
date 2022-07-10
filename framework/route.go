package framework

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

func (core *Core) GetRegisterFunc(url string, fun HandlerFunc) {
	fName := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	fEntry := handlerFuncEntry{
		funName: fName + "()",
		pattern: url,
		fun:     fun,
	}

	if err := core.routers["GET"].addRoute(url, fEntry); err != nil {
		fmt.Println()
		log.Fatal("add GET route Fatal with:", err)
	}

	core.routers["GET"].printRouteTree()
}

func (core *Core) PostRegisterFunc(url string, fun HandlerFunc) {
	fName := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	fEntry := handlerFuncEntry{
		funName: fName + "()",
		pattern: url,
		fun:     fun,
	}

	if err := core.routers["POST"].addRoute(url, fEntry); err != nil {
		fmt.Println()
		log.Fatal("add POST route Fatal with:", err)
	}

	core.routers["POST"].printRouteTree()
}

func (core *Core) getHandlerEntry(method, url string) *handlerFuncEntry {
	router, ok := core.routers[method]
	if !ok {
		return nil
	}

	return router.FindHandlerEntry(url)
}
