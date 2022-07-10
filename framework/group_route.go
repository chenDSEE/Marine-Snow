package framework

type RouteGroup interface {
	GetRegisterFunc(string, HandlerFunc)
	PostRegisterFunc(string, HandlerFunc)
	PutRegisterFunc(string, HandlerFunc)
	DeleteRegisterFunc(string, HandlerFunc)

	Group(string) RouteGroup
}

var _ RouteGroup = &prefixGroup{}

type prefixGroup struct {
	core   *Core
	parent *prefixGroup
	prefix string
}

func newPrefixGroup(core *Core, prefix string) *prefixGroup {
	return &prefixGroup{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

func (g *prefixGroup) GetRegisterFunc(uri string, fun HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	g.core.GetRegisterFunc(url, fun)
}

func (g *prefixGroup) PostRegisterFunc(uri string, fun HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	g.core.PostRegisterFunc(url, fun)
}

func (g *prefixGroup) PutRegisterFunc(uri string, fun HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	g.core.PutRegisterFunc(url, fun)
}

func (g *prefixGroup) DeleteRegisterFunc(uri string, fun HandlerFunc) {
	url := g.prefixUrl() + uri // generate full url
	g.core.DeleteRegisterFunc(url, fun)
}

func (g *prefixGroup) Group(prefix string) RouteGroup {
	group := newPrefixGroup(g.core, prefix)
	group.parent = g
	return group
}

func (g *prefixGroup) prefixUrl() string {
	if g.parent == nil {
		return g.prefix
	}

	return g.parent.prefixUrl() + g.prefix
}
