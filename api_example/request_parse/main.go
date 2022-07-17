package main

import (
	"MarineSnow/framework"
	"fmt"
	"io/ioutil"
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

/*
 * curl http://127.0.0.1:80/form-data      \
 * --data-urlencode "int=123"              \
 * --data-urlencode "int64=321"            \
 * --data-urlencode "float32=1.23"         \
 * --data-urlencode "float64=3.21"         \
 * --data-urlencode "bool=true"            \
 * --data-urlencode "stringSlice=string-1" \
 * --data-urlencode "string=string-demo"   \
 * --data-urlencode "stringSlice=string-2"
 *
 *	POST /form-data HTTP/1.1
 *	Host: 127.0.0.1
 *	User-Agent: curl/7.83.1
 *  Content-Length: 114
 *  Content-Type: application/x-www-form-urlencoded
 *  \r\n
 *  int=123&int64=321&float32=1.23&float64=3.21&bool=true&stringSlice=string-1&string=string-demo&stringSlice=string-2
 */
func formHandler(ctx *framework.Context) error {
	key := ""

	{
		key = "int"
		val, ok := ctx.FormInt(key, -1)
		fmt.Printf("FormInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.FormInt(key, -1)
		fmt.Printf("FormInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.FormInt64(key, -1)
		fmt.Printf("FormInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.FormInt64(key, -1)
		fmt.Printf("FormInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.FormFloat32(key, float32(0.0))
		fmt.Printf("FormFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.FormFloat32(key, float32(0.0))
		fmt.Printf("FormFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.FormFloat64(key, 0.0)
		fmt.Printf("FormFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.FormFloat64(key, 0.0)
		fmt.Printf("FormFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.FormBool(key, false)
		fmt.Printf("FormBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.FormBool(key, false)
		fmt.Printf("FormBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.FormString(key, "non-existed")
		fmt.Printf("FormString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.FormString(key, "non-existed")
		fmt.Printf("FormString(%s): %s, %v\n", key, non, ok)
	}
	{
		key = "stringSlice"
		val, ok := ctx.FormStringSlice(key, []string{"non-existed"})
		fmt.Printf("FormStringSlice(%s): %v, %v\n", key, val, ok)

		key = "non-stringSlice"
		non, ok := ctx.FormStringSlice(key, []string{"non-existed"})
		fmt.Printf("FormStringSlice(%s): %v, %v\n", key, non, ok)
	}

	return nil
}

// curl http://127.0.0.1:80/form-file -X POST -F 'fileName=@/root/demo.c'
// [root@LC ~]# ll demo.c
// -rw-r--r-- 1 root root 174 Jul 13 00:04 demo.c
// [root@LC ~]# cat cds.c
// #include <stdio.h>
//
// void main() {
// 		char a[2] = {0x1, 0x2};
// 		char b[2] = {0x6, 0x7};
// 		printf("%p, %p, 0x%x\n", a, b, a + 0x10);
// 		printf("0x%x, 0x%x\n", b[0], b[0 + 0x10]);
// }
// in http request:
//  Frame 6: 444 bytes on wire (3552 bits), 444 bytes captured (3552 bits)
//  Linux cooked capture v1
//  Internet Protocol Version 4, Src: 127.0.0.1, Dst: 127.0.0.1
//  Transmission Control Protocol, Src Port: 56136, Dst Port: 80, Seq: 191, Ack: 1, Len: 376
//  [2 Reassembled TCP Segments (566 bytes): #4(190), #6(376)]
//  Hypertext Transfer Protocol
//  POST /form-file HTTP/1.1\r\n
//  Host: 127.0.0.1\r\n
//  User-Agent: curl/7.83.1\r\n
//  Accept: */*\r\n
//      Content-Length: 376\r\n
//      Content-Type: multipart/form-data; boundary=------------------------ca6bb0ccadc30442\r\n
//      \r\n
//      [Full request URI: http://127.0.0.1/form-file]
//      [HTTP request 1/1]
//      [Response in frame: 8]
//      File Data: 376 bytes
//  MIME Multipart Media Encapsulation, Type: multipart/form-data, Boundary: "------------------------ca6bb0ccadc30442"
//      [Type: multipart/form-data]
//      First boundary: --------------------------ca6bb0ccadc30442\r\n
//      Encapsulated multipart part:  (application/octet-stream)
//          Content-Disposition: form-data; name="fileName"; filename="demo.c"\r\n
//          Content-Type: application/octet-stream\r\n\r\n
//          Data (174 bytes)
//      Last boundary: \r\n--------------------------ca6bb0ccadc30442--\r\n
func fileHandler(ctx *framework.Context) error {
	key := ""

	{
		key = "fileName"
		fileHeader, err := ctx.FormFile(key)
		if err == nil {
			fmt.Printf("FormFile(%s): [name: %s, size: %d]\n", key, fileHeader.Filename, fileHeader.Size)
			file, err := fileHeader.Open()
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			fmt.Printf("=== start of file ===\n")
			fmt.Printf("%s\n", string(data))
			fmt.Printf("=== end of file ===\n")

		} else {
			fmt.Printf("FormFile(%s): non-existed\n", key)
		}
	}
	{
		key = "non-file"
		fileHeader, err := ctx.FormFile(key)
		if err == nil {
			fmt.Printf("FormFile(%s): [name: %s, size: %d]\n", key, fileHeader.Filename, fileHeader.Size)
			file, err := fileHeader.Open()
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			fmt.Printf("=== start of file ===\n")
			fmt.Printf("%s\n", string(data))
			fmt.Printf("=== end of file ===\n")

		} else {
			fmt.Printf("FormFile(%s): non-existed\n", key)
		}
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
	core.PostRegisterFunc("/form-data", formHandler)
	core.PostRegisterFunc("/form-file", fileHandler)

	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
