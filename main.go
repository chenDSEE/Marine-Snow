package main

import (
	"MarineSnow/framework"
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

	/* register HTTP handler and route */
	core.GetRegisterFunc("/hello", nilHandler)
	core.GetRegisterFunc("/timeout", nilHandler)
	core.GetRegisterFunc("/timeout/demo", nilHandler)
	core.GetRegisterFunc("/hello/demo", nilHandler)
	//core.GetRegisterFunc("/hello/demo", nilHandler) // conflict with /hello/demo
	core.PostRegisterFunc("/hello/demo", nilHandler) // not conflict with GET /hello/demo

	// named parameter
	core.PostRegisterFunc("/parameter/:id", nilHandler)
	//core.PostRegisterFunc("/parameter/:name", nilHandler) // conflict with /parameter/:id
	core.PostRegisterFunc("/parameter/:id/demo", nilHandler)
	//core.PostRegisterFunc("/parameter/:id/:name", nilHandler) // conflict with /parameter/:id/demo
	core.PostRegisterFunc("/parameter/:id/:name/end", nilHandler)
	core.PostRegisterFunc("/parameter/:age/:name/new-end", nilHandler)

	// group route register
	group := core.NewRouteGroup("/group/route")
	group.GetRegisterFunc("/name", nilHandler)
	group.GetRegisterFunc("/time", nilHandler)
	group.GetRegisterFunc("/id/:name", nilHandler)
	//core.GetRegisterFunc("/group/route/name", nilHandler) // conflict with group.GetRegisterFunc("/name", nilHandler)

	core.GetRegisterFunc("/group/route/dup", nilHandler)
	//group.GetRegisterFunc("/dup", nilHandler) // conflict with core.GetRegisterFunc("/group/route/dup", nilHandler)

	// inner RouteGroup
	upGroup := core.NewRouteGroup("/up/group")
	upGroup.GetRegisterFunc("/route-1", nilHandler)
	innerGroup := upGroup.Group("/inner")
	innerGroup.GetRegisterFunc("/route-2", nilHandler)

	_ = server.ListenAndServe() // 依然借助 http.Server 来启动 http 监听、处理 connect

	fmt.Println("bye~, exit MarineSnow now.")
}
