package main

import (
	"MarineSnow/framework/gin"
	"MarineSnow/provider/demo"
	"fmt"
	"net/http"
)

// curl -i http://127.0.0.1:80/provide-and-container/must
func mustDemo(ctx *gin.Context) {
	// counter as a demo.CounterService, but not *demo.CounterServiceProvider
	counter := ctx.MustMake(demo.Key).(demo.CounterService)
	counter.Increase()
	fmt.Printf("mustDemo, [%s] counter: %d\n", counter.Name(), counter.Cnt())
}

// curl -i http://127.0.0.1:80/provide-and-container/default
func defaultDemo(ctx *gin.Context) {
	provider, err := ctx.Make(demo.Key)
	if err != nil {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	counter, ok := provider.(demo.CounterService)
	if !ok {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	counter.Increase()
	fmt.Printf("defaultDemo, [%s] counter: %d\n", counter.Name(), counter.Cnt())
}

// curl -i http://127.0.0.1:80/provide-and-container/MakeNew
// call this URL more tha twice, but counter still 1
func makeNewDemo(ctx *gin.Context) {
	provider, err := ctx.MakeNew(demo.Key, "new-params")
	if err != nil {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	counter, ok := provider.(demo.CounterService)
	if !ok {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	counter.Increase()
	fmt.Printf("makeNewDemo, [%s] counter: %d\n", counter.Name(), counter.Cnt())
}

// curl -i http://127.0.0.1:80/provide-and-container/not-register
func notRegisterDemo(ctx *gin.Context) {
	fmt.Println("notRegisterDemo:")
	fmt.Println("not-register-key:", ctx.IsRegistered("not-register-key"))
	fmt.Printf("%s: %v\n", demo.Key, ctx.IsRegistered(demo.Key))
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := gin.New()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	// register CounterServiceProvider into Container
	_ = core.Register(&demo.CounterServiceProvider{})

	/* register HTTP handler and route */
	// request demo
	core.GET("/provide-and-container/must", mustDemo)
	core.GET("/provide-and-container/default", defaultDemo)
	core.GET("/provide-and-container/MakeNew", makeNewDemo)
	core.GET("/provide-and-container/not-register", notRegisterDemo)

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
