package framework

import (
	"context"
	"net/http"
)

type Context struct {
	Rsp http.ResponseWriter
	Req *http.Request

	Ctx context.Context
}
