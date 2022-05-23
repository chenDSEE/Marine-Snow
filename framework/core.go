package framework

import (
	"fmt"
	"net/http"
)

type Core struct {
	name string
}

func NewCore() *Core {
	return &Core{
		name: "MarineSnow Core",
	}
}

func (core *Core) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	fmt.Println(core.name, "demo")
}
