package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello MarineSnow")
	server := &http.Server{
		Addr:    "127.0.0.1:80",
		Handler: framework.NewCore(),
	}

	server.ListenAndServe()
}
