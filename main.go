package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

func helloHandler(ctx *framework.Context) error {
	fmt.Println("in hello handler func:", ctx.Request().URL.String())
	return nil
}

func timedemoHandler(ctx *framework.Context) error {
	fmt.Println("in timedemoHandler func:", ctx.Request().URL.String())
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
	core.RegisterHandlerFunc("/hello", helloHandler)

	_ = server.ListenAndServe() // 依然借助 http.Server 来启动 http 监听、处理 connect

	fmt.Println("bye~, exit MarineSnow now.")
}
