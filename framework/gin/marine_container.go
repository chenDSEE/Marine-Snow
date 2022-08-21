package gin

import (
	"MarineSnow/framework"
)

// for gin.Engine
func (engine *Engine) Register(provider framework.ServiceProvider) error {
	return engine.container.Register(provider)
}

func (engine *Engine) IsRegistered(key string) bool {
	return engine.container.IsRegistered(key)
}

// for gin.Context
func (ctx *Context) Register(provider framework.ServiceProvider) error {
	return ctx.container.Register(provider)
}

func (ctx *Context) IsRegistered(key string) bool {
	return ctx.container.IsRegistered(key)
}

func (ctx *Context) Make(key string) (interface{}, error) {
	return ctx.container.Make(key)
}

func (ctx *Context) MustMake(key string) interface{} {
	return ctx.container.MustMake(key)
}

func (ctx *Context) MakeNew(key string, params ...interface{}) (interface{}, error) {
	return ctx.container.MakeNew(key, params...)
}
