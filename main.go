package main

import (
	"MarineSnow/framework"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// curl http://127.0.0.1:80/basic-info
// curl http://127.0.0.1:80/basic-info -H 'X-Real-Ip: 1.1.1.1'
// curl http://127.0.0.1:80/basic-info -H 'X-Forwarded-For: 2.2.2.2'
func basicInfoHandler(ctx *framework.Context) error {
	fmt.Printf("[%s --> %s]: %s %s\n", ctx.ClientIp(), ctx.Host(), ctx.Method(), ctx.URI())
	return nil
}

// curl http://127.0.0.1:80/header
func headerHandler(ctx *framework.Context) error {
	key := "User-Agent"
	agent, ok := ctx.Header(key)
	fmt.Printf("Header(%s): %s, %v\n", key, agent, ok)

	key = "X-Forwarded-For" // non existed
	non, ok := ctx.Header(key)
	fmt.Printf("Header(%s): %s, %v\n", key, non, ok)

	// dump all Header
	fmt.Println("=== dump all Header ===")
	headers := ctx.Headers()
	for header, val := range headers {
		fmt.Printf("%s: %s\n", header, val)
	}

	return nil
}

// curl http://127.0.0.1:80/cookie --cookie 'key-1=value-1' --cookie 'key-2=value-2' --cookie 'key-1=value-3'
func cookieHandler(ctx *framework.Context) error {
	key := "key-2"
	cookie, ok := ctx.Cookie(key)
	fmt.Printf("Cookie(%s): %s, %v\n", key, cookie, ok)

	key = "non-key" // non existed
	cookie, ok = ctx.Cookie(key)
	fmt.Printf("Cookie(%s): %s, %v\n", key, cookie, ok)

	// dump all cookies
	fmt.Println("=== dump all cookies ===")
	cookies := ctx.Cookies()
	for key, val := range cookies {
		fmt.Printf("%s(num:%d): %v\n", key, len(val), val)
	}

	return nil
}

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

// Issue by:
// curl -X POST \
// http://127.0.0.1:80/json \
// -H 'content-type: application/json' \
// -d '{"name":"name-string","num":123,"now":"2022-07-09T23:00:00Z","data":["string-1","string-2"]}'
//
// In http request:
// 	POST /json HTTP/1.1
//	Host: 127.0.0.1
//	User-Agent: curl/7.83.1
//	Accept: */*
//	content-type: application/json
//	Content-Length: 92
//
//	{"name":"name-string","num":123,"now":"2022-07-09T23:00:00Z","data":["string-1","string-2"]}
func jsonHandler(ctx *framework.Context) error {
	type jsonObj struct {
		Name string    `json:"name"`
		Num  int       `json:"num"`
		Data []string  `json:"data"`
		Now  time.Time `json:"now"`
	}

	var obj jsonObj
	err := ctx.BindJSON(&obj)
	if err != nil {
		return err
	}

	fmt.Printf("JSON data: obj%v\n", obj)
	fmt.Printf("raw JSON data: %s\n", ctx.Request().Body)
	return nil
}

// Issue by:
//  curl -X POST \
//  http://127.0.0.1:80/xml \
//  -H 'content-type: application/xml' \
//  -d '<?xml version="1.0" encoding="UTF-8"?><entry><name>name-string</name><num>123</num><now>2022-07-09T23:00:00Z</now><data>string-1</data><data>string-2</data></entry>'
//
// In http request:
//  POST /xml HTTP/1.1
//  Host: 127.0.0.1
//  User-Agent: curl/7.83.1
//  Accept: */*
//  content-type: application/xml
//  Content-Length: 164
//
//  <?xml version="1.0" encoding="UTF-8"?><entry><name>name-string</name><num>123</num><now>2022-07-09T23:00:00Z</now><data>string-1</data><data>string-2</data></entry>
func xmlHandler(ctx *framework.Context) error {
	type xmlObj struct {
		XMLName xml.Name  `xml:"entry"`
		Name    string    `xml:"name"`
		Num     int       `xml:"num"`
		Data    []string  `xml:"data"`
		Now     time.Time `xml:"now"`
	}

	var obj xmlObj
	err := ctx.BindXML(&obj)
	if err != nil {
		return err
	}

	fmt.Printf("XML data: obj%v\n", obj)
	fmt.Printf("raw XML data: %s\n", ctx.Request().Body)
	return nil
}

// curl -X POST http://127.0.0.1:80/raw -d 'raw-data-string'
func rawBodyHandler(ctx *framework.Context) error {
	data, err := ctx.RawBody()
	if err != nil {
		return err
	}

	fmt.Printf("raw data: obj%v\n", data)
	fmt.Printf("raw data: %s\n", ctx.Request().Body)

	//data, err = ctx.RawBody()
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Printf("raw data: obj%v\n", data)
	//fmt.Printf("raw data: %s\n", ctx.Request().Body)
	return nil
}

/* response demo */
// curl -i http://127.0.0.1/basic-response
func basicResponseHandler(ctx *framework.Context) error {
	ctx.AddHeader("Allow", http.MethodPost)
	ctx.AddHeader("date", "Sun, 24 Jul 2022 05:51:10 GMT")
	ctx.DelHeader("date")

	ctx.SetCookie("cookie-key", "cookie-value", 10, "", "", false, true)

	ctx.SetStatus(http.StatusAccepted)
	return nil
}

// curl -i http://127.0.0.1/text/no-status
func textBodyNoStatusHandler(ctx *framework.Context) error {
	ctx.Text("Text data:[%s]\n", "demo-string1")
	return nil
}

// curl -i http://127.0.0.1/text/with-status
func textBodyWithStatusHandler(ctx *framework.Context) error {
	ctx.SetStatus(http.StatusMultiStatus) // make on effect to Header "Content-Type"
	ctx.Text("Text data[%d]:[%s]\n", http.StatusMultiStatus, "demo-string1")

	return nil
}

// curl -i http://127.0.0.1/raw-data
func rawDataHandler(ctx *framework.Context) error {
	// any Context-Type as you want
	//ctx.AddHeader("Content-Type", "application/text")
	//ctx.AddHeader("Content-Type", "application/json")
	//ctx.AddHeader("Content-Type", "application/xml")
	//ctx.AddHeader("Content-Type", "text/plain")
	ctx.AddHeader("Content-Type", "application/html")

	ctx.RawData([]byte("byte-demo-for-raw-data\n"))
	return nil
}

// curl -i http://127.0.0.1/json-data
func jsonDataHandler(ctx *framework.Context) error {
	type jsonObj struct {
		Name string    `json:"name"`
		Num  int       `json:"num"`
		Data []string  `json:"data"`
		Now  time.Time `json:"now"`
	}

	ctx.JSON(&jsonObj{
		Name: "name-string",
		Num:  20,
		Data: []string{"string-1", "string2"},
		Now:  time.Now(),
	})
	return nil
}

// curl -i http://127.0.0.1/xml-data
func xmlDataHandler(ctx *framework.Context) error {
	type xmlObj struct {
		XMLName xml.Name  `xml:"entry"`
		Name    string    `xml:"name"`
		Num     int       `xml:"num"`
		Data    []string  `xml:"data"`
		Now     time.Time `xml:"now"`
	}

	ctx.XML(&xmlObj{
		Name: "name-string",
		Num:  30,
		Data: []string{"string-1", "string-2"},
		Now:  time.Now(),
	})
	return nil
}

// curl -i http://127.0.0.1/redirect
func redirectHandler(ctx *framework.Context) error {
	ctx.Redirect("/www.redirect.com/new/path")
	return nil
}

// curl -i http://127.0.0.1/html/files
func htmlFilesHandler(ctx *framework.Context) error {
	type HtmlObj struct {
		Name string
	}

	e := HtmlObj{Name: "Marine Snow"}
	ctx.HtmlFiles(e, "example.html")
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
	// request demo
	core.GetRegisterFunc("/basic-info", basicInfoHandler)
	core.GetRegisterFunc("/header", headerHandler)
	core.GetRegisterFunc("/cookie", cookieHandler)
	core.GetRegisterFunc("/query", queryHandler)
	core.GetRegisterFunc("/int/:int/:int64/float/:float32/:float64/bool/:bool/string/:string/end", paramHandler)
	core.PostRegisterFunc("/form-data", formHandler)
	core.PostRegisterFunc("/form-file", fileHandler)
	core.PostRegisterFunc("/json", jsonHandler)
	core.PostRegisterFunc("/xml", xmlHandler)
	core.PostRegisterFunc("/raw", rawBodyHandler)

	// response demo
	core.GetRegisterFunc("/basic-response", basicResponseHandler)
	core.GetRegisterFunc("/text/no-status", textBodyNoStatusHandler)
	core.GetRegisterFunc("/text/with-status", textBodyWithStatusHandler)
	core.GetRegisterFunc("/json-data", jsonDataHandler)
	core.GetRegisterFunc("/xml-data", xmlDataHandler)
	core.GetRegisterFunc("/redirect", redirectHandler)
	core.GetRegisterFunc("/html/files", htmlFilesHandler)

	core.GetRegisterFunc("/raw-data", rawDataHandler)

	core.DumpRoutes()

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
