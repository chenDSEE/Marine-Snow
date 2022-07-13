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

// func (core *Core) getHandlerEntryList(method, url string) ([]handlerFuncEntry, string) {
// 	router, ok := core.routers[method]
// 	if !ok {
// 		return nil, ""
// 	}

// 	return router.FindHandlerEntryList(url)
// }

func (core *Core) matchRoute(method, path string) (*routeEntry, bool) {
	router, ok := core.routers[method]
	if !ok {
		return nil, false
	}

	matchNode := router.root.matchNode(path)
	if matchNode == nil {
		return nil, false
	}

	entry := &routeEntry{
		method:     method,
		path:       path,
		endNode:    matchNode,
		paramTable: map[string]string{},
	}

	return entry, true
}

// TODO: interface this struct
// not safe for concurrency
type routeEntry struct {
	method     string
	path       string
	endNode    *node
	paramTable map[string]string
	// router *rTree   no need now
}

func (entry *routeEntry) getHandlerEntryList() []handlerFuncEntry {
	if entry.endNode == nil {
		return nil
	}

	return entry.endNode.handlerEntryList
}

func (entry *routeEntry) getParamTable() map[string]string {
	if len(entry.paramTable) > 0 {
		// only parse once
		return entry.paramTable
	}

	if entry.endNode == nil {
		return entry.paramTable
	}

	entry.paramTable = entry.endNode.parseUrlParameters(entry.path)

	return entry.paramTable
}

func (entry *routeEntry) getPattern() string {
	if entry.endNode == nil {
		return ""
	}

	return entry.endNode.pattern
}

func (entry *routeEntry) getMethod() string {
	return entry.method
}

func (entry *routeEntry) getPath() string {
	return entry.path
}
