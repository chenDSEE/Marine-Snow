package main

import (
	"MarineSnow/framework"
	"MarineSnow/framework/middleware"
	"context"
	"fmt"
	"net/http"
	"time"
)

func helloHandler(ctx *framework.Context) error {
	return nil
}

func nilHandler(ctx *framework.Context) error {
	return nil
}

type info struct {
	name string
	data string
}

func timedemoHandler(ctx *framework.Context) error {
	infoChan := make(chan *info, 1)

	timeoutCtx, cancel := context.WithTimeout(ctx.BaseContext(), time.Duration(2*time.Second))
	defer cancel()

	go func() {
		fmt.Printf("RPC start  -->\n")

		time.Sleep(1 * time.Second)
		infoChan <- &info{name: "name", data: "data"}

		fmt.Printf("RPC end   <--\n")
	}()

	select {
	case info := <-infoChan:
		fmt.Printf("get information:[name:%s, data:%s]\n", info.name, info.data)
	case <-timeoutCtx.Done():
		fmt.Printf("timeout and cancel this RPC\n")
	}

	fmt.Printf("timedemoHandler() finish and exit\n")
	return nil
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := framework.NewCore()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	core.AppendDefaultMiddleware(
		middleware.Recovery(),
		middleware.Cost())

	/* register HTTP handler and route */
	core.GetRegisterFunc("/hello",
		middleware.Example1(),
		middleware.Example2(),
		middleware.WrapNext(nilHandler),
		nilHandler,
		middleware.Example3())

	group := core.NewRouteGroup("/group")
	group.AppendDefaultMiddleware(middleware.Example1())
	group.GetRegisterFunc("/demo",
		middleware.Example2(),
		middleware.WrapNext(nilHandler),
		nilHandler,
		middleware.Example3())

	core.DumpRoutes()

	_ = server.ListenAndServe() // 依然借助 http.Server 来启动 http 监听、处理 connect

	fmt.Println("bye~, exit Marine-Snow now.")
}
