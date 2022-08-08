package main

import (
	"MarineSnow/framework/gin"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// curl -i http://127.0.0.1:80/basic-info
// curl -i http://127.0.0.1:80/basic-info -H 'X-Real-Ip: 1.1.1.1'
// curl -i http://127.0.0.1:80/basic-info -H 'X-Forwarded-For: 2.2.2.2'
func basicInfoHandler(ctx *gin.Context) {
	fmt.Printf("result for MarineSnow API [%s --> %s]: %s %s\n", ctx.ClientIP(), ctx.Host(), ctx.Method(), ctx.URI())
	fmt.Printf("result for gin API [%s --> %s]: %s %s\n", ctx.ClientIP(), ctx.Host(), ctx.Method(), ctx.FullPath())
}

// curl -v http://127.0.0.1:80/header
func headerHandler(ctx *gin.Context) {
	key := "User-Agent"
	agent := ctx.GetHeader(key)
	fmt.Printf("Header(%s): %s\n", key, agent)

	key = "X-Forwarded-For" // non existed
	non := ctx.GetHeader(key)
	fmt.Printf("Header(%s): %s\n", key, non)

	// dump all Header
	fmt.Println("=== dump all Header ===")
	headers := ctx.GetHeaders()
	for header, val := range headers {
		fmt.Printf("%s: %s\n", header, val)
	}
}

// curl -v http://127.0.0.1:80/cookie --cookie 'key-1=value-1' --cookie 'key-2=value-2' --cookie 'key-1=value-3'
func cookieHandler(ctx *gin.Context) {
	key := "key-2"
	cookie, err := ctx.Cookie(key)
	fmt.Printf("Cookie(%s): %s, error: %v\n", key, cookie, err)

	key = "non-key" // non existed
	cookie, err = ctx.Cookie(key)
	fmt.Printf("Cookie(%s): %s, error: %v\n", key, cookie, err)

	// dump all cookies
	fmt.Println("=== dump all cookies ===")
	cookies := ctx.Cookies()
	for key, val := range cookies {
		fmt.Printf("%s(num:%d): %v\n", key, len(val), val)
	}
}

// curl -v "http://127.0.0.1:80/query?int=123&int64=321&float32=1.23&float64=3.21&bool=True&string=demo-string&stringSlice=demo-stringSlice-1&stringSlice=demo-stringSlice-2"
func queryHandler(ctx *gin.Context) {
	key := ""

	{
		key = "int"
		val, ok := ctx.DefaultQueryInt(key, -1)
		fmt.Printf("DefaultQueryInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.DefaultQueryInt(key, -1)
		fmt.Printf("DefaultQueryInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.DefaultQueryInt64(key, -1)
		fmt.Printf("DefaultQueryInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.DefaultQueryInt64(key, -1)
		fmt.Printf("DefaultQueryInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.DefaultQueryFloat32(key, float32(0.0))
		fmt.Printf("DefaultQueryFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.DefaultQueryFloat32(key, float32(0.0))
		fmt.Printf("DefaultQueryFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.DefaultQueryFloat64(key, 0.0)
		fmt.Printf("DefaultQueryFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.DefaultQueryFloat64(key, 0.0)
		fmt.Printf("DefaultQueryFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.DefaultQueryBool(key, false)
		fmt.Printf("DefaultQueryBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.DefaultQueryBool(key, false)
		fmt.Printf("DefaultQueryBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.DefaultQueryString(key, "non-existed")
		fmt.Printf("DefaultQueryString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.DefaultQueryString(key, "non-existed")
		fmt.Printf("DefaultQueryString(%s): %s, %v\n", key, non, ok)
	}
	{
		key = "stringSlice"
		val, ok := ctx.DefaultQueryStringSlice(key, []string{"non-existed"})
		fmt.Printf("DefaultQueryStringSlice(%s): %v, %v\n", key, val, ok)

		key = "non-stringSlice"
		non, ok := ctx.DefaultQueryStringSlice(key, []string{"non-existed"})
		fmt.Printf("DefaultQueryStringSlice(%s): %v, %v\n", key, non, ok)
	}
}

// curl -v http://127.0.0.1:80/int/123/321/float/1.23/3.21/bool/true/string/demo-string/end
// /int/:int/:int64/float/:float32/:float64/bool/:bool/string/:string/end
func paramHandler(ctx *gin.Context) {
	key := ""

	{
		key = "int"
		val, ok := ctx.DefaultParamInt(key, -1)
		fmt.Printf("DefaultParamInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.DefaultParamInt(key, -1)
		fmt.Printf("DefaultParamInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.DefaultParamInt64(key, -1)
		fmt.Printf("DefaultParamInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.DefaultParamInt64(key, -1)
		fmt.Printf("DefaultParamInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.DefaultParamFloat32(key, float32(0.0))
		fmt.Printf("DefaultParamFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.DefaultParamFloat32(key, float32(0.0))
		fmt.Printf("DefaultParamFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.DefaultParamFloat64(key, 0.0)
		fmt.Printf("DefaultParamFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.DefaultParamFloat64(key, 0.0)
		fmt.Printf("DefaultParamFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.DefaultParamBool(key, false)
		fmt.Printf("DefaultParamBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.DefaultParamBool(key, false)
		fmt.Printf("DefaultParamBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.DefaultParamString(key, "non-existed")
		fmt.Printf("DefaultParamString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.DefaultParamString(key, "non-existed")
		fmt.Printf("DefaultParamString(%s): %s, %v\n", key, non, ok)
	}
}

/*
* curl -v http://127.0.0.1:80/form-data      \
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
func formHandler(ctx *gin.Context) {
	key := ""

	{
		key = "int"
		val, ok := ctx.DefaultFormInt(key, -1)
		fmt.Printf("DefaultFormInt(%s): %d, %v\n", key, val, ok)

		key = "non-int"
		non, ok := ctx.DefaultFormInt(key, -1)
		fmt.Printf("DefaultFormInt(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "int64"
		val, ok := ctx.DefaultFormInt64(key, -1)
		fmt.Printf("DefaultFormInt64(%s): %d, %v\n", key, val, ok)

		key = "non-int64"
		non, ok := ctx.DefaultFormInt64(key, -1)
		fmt.Printf("DefaultFormInt64(%s): %d, %v\n", key, non, ok)
	}
	{
		key = "float32"
		val, ok := ctx.DefaultFormFloat32(key, float32(0.0))
		fmt.Printf("DefaultFormFloat32(%s): %f, %v\n", key, val, ok)

		key = "non-float32"
		non, ok := ctx.DefaultFormFloat32(key, float32(0.0))
		fmt.Printf("DefaultFormFloat32(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "float64"
		val, ok := ctx.DefaultFormFloat64(key, 0.0)
		fmt.Printf("DefaultFormFloat64(%s): %f, %v\n", key, val, ok)

		key = "non-float64"
		non, ok := ctx.DefaultFormFloat64(key, 0.0)
		fmt.Printf("DefaultFormFloat64(%s): %f, %v\n", key, non, ok)
	}
	{
		key = "bool"
		val, ok := ctx.DefaultFormBool(key, false)
		fmt.Printf("DefaultFormBool(%s): %v, %v\n", key, val, ok)

		key = "non-bool"
		non, ok := ctx.DefaultFormBool(key, false)
		fmt.Printf("DefaultFormBool(%s): %v, %v\n", key, non, ok)
	}
	{
		key = "string"
		val, ok := ctx.DefaultFormString(key, "non-existed")
		fmt.Printf("DefaultFormString(%s): %s, %v\n", key, val, ok)

		key = "non-string"
		non, ok := ctx.DefaultFormString(key, "non-existed")
		fmt.Printf("DefaultFormString(%s): %s, %v\n", key, non, ok)
	}
	{
		key = "stringSlice"
		val, ok := ctx.DefaultFormStringSlice(key, []string{"non-existed"})
		fmt.Printf("DefaultFormStringSlice(%s): %v, %v\n", key, val, ok)

		key = "non-stringSlice"
		non, ok := ctx.DefaultFormStringSlice(key, []string{"non-existed"})
		fmt.Printf("DefaultFormStringSlice(%s): %v, %v\n", key, non, ok)
	}
}

//
// curl -v http://127.0.0.1:80/form-file -X POST -F 'fileName=@/root/demo.c'
// [root@LC ~]# ll demo.c
// -rw-r--r-- 1 root root 174 Jul 13 00:04 demo.c
// [root@LC ~]# cat demo.c
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
func fileHandler(ctx *gin.Context) {
	key := ""

	{
		key = "fileName"
		fileHeader, err := ctx.FormFile(key)
		if err == nil {
			fmt.Printf("FormFile(%s): [name: %s, size: %d]\n", key, fileHeader.Filename, fileHeader.Size)
			file, err := fileHeader.Open()
			if err != nil {
				return
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return
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
				return
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return
			}
			fmt.Printf("=== start of file ===\n")
			fmt.Printf("%s\n", string(data))
			fmt.Printf("=== end of file ===\n")

		} else {
			fmt.Printf("FormFile(%s): non-existed\n", key)
		}
	}
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
func jsonHandler(ctx *gin.Context) {
	type jsonObj struct {
		Name string    `json:"name"`
		Num  int       `json:"num"`
		Data []string  `json:"data"`
		Now  time.Time `json:"now"`
	}

	var obj jsonObj
	err := ctx.BindJSON(&obj)
	if err != nil {
		return
	}

	fmt.Printf("JSON data: obj%+v\n", obj)
}

// Issue by:
//  curl -v -X POST \
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
func xmlHandler(ctx *gin.Context) {
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
		return
	}

	fmt.Printf("XML data: obj%+v\n", obj)
}

/* response demo */
// curl -i http://127.0.0.1/basic-response
func basicResponseHandler(ctx *gin.Context) {
	ctx.IAddHeader("Allow", http.MethodPost)
	ctx.IAddHeader("User-Agent", "MarineSnow")
	ctx.IDelHeader("User-Agent")

	ctx.ISetCookie("cookie-key", "cookie-value", 10, "", "", false, true)

	ctx.ISetStatus(http.StatusAccepted)
}

// curl -i http://127.0.0.1/text/no-status
func textBodyNoStatusHandler(ctx *gin.Context) {
	ctx.IText("Text data:[%s]\n", "demo-string1")
}

// curl -i http://127.0.0.1/text/with-status
func textBodyWithStatusHandler(ctx *gin.Context) {
	ctx.ISetStatus(http.StatusMultiStatus) // make on effect to Header "Content-Type"
	ctx.IText("Text data[%d]:[%s]\n", http.StatusMultiStatus, "demo-string1")
}

// curl -i http://127.0.0.1/raw-data
func rawDataHandler(ctx *gin.Context) {
	// any Context-Type as you want
	//ctx.AddHeader("Content-Type", "application/text")
	//ctx.AddHeader("Content-Type", "application/json")
	//ctx.AddHeader("Content-Type", "application/xml")
	//ctx.AddHeader("Content-Type", "text/plain")
	ctx.IAddHeader("Content-Type", "application/html")

	ctx.IRawData([]byte("byte-demo-for-raw-data\n"))
}

// curl -i http://127.0.0.1/json-data
func jsonDataHandler(ctx *gin.Context) {
	type jsonObj struct {
		Name string    `json:"name"`
		Num  int       `json:"num"`
		Data []string  `json:"data"`
		Now  time.Time `json:"now"`
	}

	ctx.IJSON(&jsonObj{
		Name: "name-string",
		Num:  20,
		Data: []string{"string-1", "string2"},
		Now:  time.Now(),
	})
}

// curl -i http://127.0.0.1/xml-data
func xmlDataHandler(ctx *gin.Context) {
	type xmlObj struct {
		XMLName xml.Name  `xml:"entry"`
		Name    string    `xml:"name"`
		Num     int       `xml:"num"`
		Data    []string  `xml:"data"`
		Now     time.Time `xml:"now"`
	}

	ctx.IXML(&xmlObj{
		Name: "name-string",
		Num:  30,
		Data: []string{"string-1", "string-2"},
		Now:  time.Now(),
	})
}

// curl -i http://127.0.0.1/redirect
func redirectHandler(ctx *gin.Context) {
	ctx.IRedirect("/www.redirect.com/new/path")
}

// curl -i http://127.0.0.1/html/files
func htmlFilesHandler(ctx *gin.Context) {
	type HtmlObj struct {
		Name string
	}

	e := HtmlObj{Name: "Marine Snow"}
	ctx.IHtmlFiles(e, "example.html")
}

// curl -i https://127.0.0.1:80/jsonp?callbackKey=callbackFunction
func jsonpHandler(ctx *gin.Context) {
	type jsonObj struct {
		Name string    `json:"name"`
		Num  int       `json:"num"`
		Data []string  `json:"data"`
		Now  time.Time `json:"now"`
	}

	ctx.IJsonp("callbackKey", &jsonObj{
		Name: "name-string",
		Num:  20,
		Data: []string{"string-1", "string2"},
		Now:  time.Now(),
	})
}

const SERVER_ADDR = "127.0.0.1:80"

func main() {
	fmt.Printf("welcome to MarineSnow, start now. Listen on [%s]\n", SERVER_ADDR)
	core := gin.New()
	server := &http.Server{
		Addr:    SERVER_ADDR,
		Handler: core,
	}

	/* register HTTP handler and route */
	// request demo
	core.GET("/basic-info", basicInfoHandler)
	core.GET("/header", headerHandler)
	core.GET("/cookie", cookieHandler)
	core.GET("/query", queryHandler)
	core.GET("/int/:int/:int64/float/:float32/:float64/bool/:bool/string/:string/end", paramHandler)
	core.POST("/form-data", formHandler)
	core.POST("/form-file", fileHandler)
	core.POST("/json", jsonHandler)
	core.POST("/xml", xmlHandler)
	//core.POST("/raw", rawBodyHandler)

	// response demo
	core.GET("/basic-response", basicResponseHandler)
	core.GET("/text/no-status", textBodyNoStatusHandler)
	core.GET("/text/with-status", textBodyWithStatusHandler)
	core.GET("/json-data", jsonDataHandler)
	core.GET("/xml-data", xmlDataHandler)
	core.GET("/redirect", redirectHandler)
	core.GET("/html/files", htmlFilesHandler)
	core.GET("/jsonp", jsonpHandler)
	core.GET("/raw-data", rawDataHandler)

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit Marine-Snow now.")
}
