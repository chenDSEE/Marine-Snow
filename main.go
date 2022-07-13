package main

import (
	"MarineSnow/framework"
	"fmt"
	"net/http"
)

// curl "http://127.0.0.1:80/query?int=123&int64=321&float32=1.23&float64=3.21&bool=True&string=demo-string&stringSlice=demo-stringSlice-1&stringSlice=demo-stringSlice-2"
func queryHandler(ctx *framework.Context) error {
	key := ""

	{
		key = "int"
		val, ok := ctx.QueryInt(key, -1)
		fmt.Printf("QueryInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.QueryInt(key, -1)
		fmt.Printf("QueryInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.QueryInt64(key, -1)
		fmt.Printf("QueryInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.QueryInt64(key, -1)
		fmt.Printf("QueryInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.QueryFloat32(key, float32(0.0))
		fmt.Printf("QueryFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.QueryFloat32(key, float32(0.0))
		fmt.Printf("QueryFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.QueryFloat64(key, 0.0)
		fmt.Printf("QueryFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.QueryFloat64(key, 0.0)
		fmt.Printf("QueryFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.QueryBool(key, false)
		fmt.Printf("QueryBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.QueryBool(key, false)
		fmt.Printf("QueryBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.QueryString(key, "non-existed")
		fmt.Printf("QueryString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.QueryString(key, "non-existed")
		fmt.Printf("QueryString(%s): %s, %v\n", key, non, ok)
	}
	{
		key = "stringSlice"
		val, ok := ctx.QueryStringSlice(key, []string{"non-existed"})
		fmt.Printf("QueryStringSlice(%s): %v, %v\n", key, val, ok)

		key = "non-stringSlice"
		non, ok := ctx.QueryStringSlice(key, []string{"non-existed"})
		fmt.Printf("QueryStringSlice(%s): %v, %v\n", key, non, ok)
	}

	return nil
}

// curl http://127.0.0.1:80/int/123/321/float/1.23/3.21/bool/true/string/demo-string/end
// /int/:int/:int64/float/:float32/:float64/bool/:bool/string/:string/end
func paramHandler(ctx *framework.Context) error {
	key := ""

	{
		key = "int"
		val, ok := ctx.ParamInt(key, -1)
		fmt.Printf("ParamInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.ParamInt(key, -1)
		fmt.Printf("ParamInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.ParamInt64(key, -1)
		fmt.Printf("ParamInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.ParamInt64(key, -1)
		fmt.Printf("ParamInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.ParamFloat32(key, float32(0.0))
		fmt.Printf("ParamFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.ParamFloat32(key, float32(0.0))
		fmt.Printf("ParamFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.ParamFloat64(key, 0.0)
		fmt.Printf("ParamFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.ParamFloat64(key, 0.0)
		fmt.Printf("ParamFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.ParamBool(key, false)
		fmt.Printf("ParamBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.ParamBool(key, false)
		fmt.Printf("ParamBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.ParamString(key, "non-existed")
		fmt.Printf("ParamString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.ParamString(key, "non-existed")
		fmt.Printf("ParamString(%s): %s, %v\n", key, non, ok)
	}

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
	core.GetRegisterFunc("/query", queryHandler)
	core.GetRegisterFunc("/int/:int/:int64/float/:float32/:float64/bool/:bool/string/:string/end", paramHandler)

	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
