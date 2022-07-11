package framework

type RouteGroup interface {
	GetRegisterFunc(string, ...HandlerFunc)
	PostRegisterFunc(string, ...HandlerFunc)
	PutRegisterFunc(string, ...HandlerFunc)
	DeleteRegisterFunc(string, ...HandlerFunc)

	Group(string) RouteGroup

	AppendDefaultMiddleware(middlewareFun ...HandlerFunc)
}

var _ RouteGroup = &prefixGroup{}

type prefixGroup struct {
	core        *Core
	parent      *prefixGroup
	prefix      string
	middlewares []HandlerFunc
}

func newPrefixGroup(core *Core, prefix string) *prefixGroup {
	return &prefixGroup{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

func (g *prefixGroup) GetRegisterFunc(uri string, fun ...HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	handlerList := make([]HandlerFunc, 0, len(fun)+len(g.middlewares))
	handlerList = append(g.middlewares, fun...)
	g.core.GetRegisterFunc(url, handlerList...)
}

func (g *prefixGroup) PostRegisterFunc(uri string, fun ...HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	handlerList := make([]HandlerFunc, 0, len(fun)+len(g.middlewares))
	handlerList = append(g.middlewares, fun...)
	g.core.PostRegisterFunc(url, handlerList...)
}

func (g *prefixGroup) PutRegisterFunc(uri string, fun ...HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	handlerList := make([]HandlerFunc, 0, len(fun)+len(g.middlewares))
	handlerList = append(g.middlewares, fun...)
	g.core.PutRegisterFunc(url, handlerList...)
}

func (g *prefixGroup) DeleteRegisterFunc(uri string, fun ...HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	handlerList := make([]HandlerFunc, 0, len(fun)+len(g.middlewares))
	handlerList = append(g.middlewares, fun...)
	g.core.DeleteRegisterFunc(url, handlerList...)
}

func (g *prefixGroup) Group(prefix string) RouteGroup {
	group := newPrefixGroup(g.core, prefix)
	group.parent = g
	return group
}

func (g *prefixGroup) AppendDefaultMiddleware(middlewareFun ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewareFun...)
}

func (g *prefixGroup) prefixUrl() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.prefixUrl() + g.prefix
}
