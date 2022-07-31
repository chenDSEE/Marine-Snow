package main

import (
	"MarineSnow/framework"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// curl -i http://127.0.0.1:80/long-response
func longRespHandler(ctx *framework.Context) error {
	fmt.Printf("[%s --> %s]: %s %s, long request start\n", ctx.ClientIp(), ctx.Host(), ctx.Method(), ctx.URI())
	time.Sleep(time.Second * 12)
	fmt.Printf("After 5s, long reqest done, reply response\n")
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

	/* register signal handle */
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	/* register HTTP handler and route */
	server.RegisterOnShutdown(func() {
		fmt.Printf("HTTP server received shutdown command, start to shutdown\n")
	})

	core.GetRegisterFunc("/long-response", longRespHandler)
	core.DumpRoutes()

	serverChan := make(chan error, 1)
	go func() {
		serverChan <- server.ListenAndServe()
	}()

	select {
	case sig := <-sigChan:
		fmt.Printf("Receive signal:[%s], start to shutdown server\n", sig.String())
		shutdownTimeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := server.Shutdown(shutdownTimeCtx); err != nil {
			fmt.Printf("Shutdown server fail:[%s]\n", err.Error())
			break
		}

		fmt.Printf("HTTP connection shutdown, server gracefully exit success.\n")
	case err := <-serverChan:
		// http server already exit, no need to shut down anything
		fmt.Printf("Receive error from http server:[%s]\n", err.Error())
	}

	fmt.Println("bye~, exit Marine-Snow now.")
}
