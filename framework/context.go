package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	/* request part */
	req *http.Request

	ctx context.Context

	hEntryIndex int
	hEntryList  []handlerFuncEntry

	// cookie
	cookieParse sync.Once
	cookieTable map[string][]string

	// query data in URL
	queryAll   sync.Once
	queryTable map[string][]string

	// param data in URL
	paramTable map[string]string

	/* response part */
	rsp http.ResponseWriter

	// only wrote status when write body
	// without body, and SetStatus(), will write 200 by default
	statusCode  int
	wroteStatus bool
}

func NewContext(rsp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		/* request part */
		req:         req,
		ctx:         req.Context(),
		hEntryIndex: -1,
		hEntryList:  nil,
		cookieTable: map[string][]string{},
		queryTable:  map[string][]string{},
		paramTable:  map[string]string{},

		/* response part */
		rsp:        rsp,
		statusCode: http.StatusOK, // default
	}
}

var _ context.Context = &Context{}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

// method to access internal variable.
// TODO: wrap those method better
// TODO: remove below method
func (ctx *Context) BaseContext() context.Context {
	return ctx.req.Context()
}

func (ctx *Context) Request() *http.Request {
	return ctx.req
}

func (ctx *Context) ResponseWriter() http.ResponseWriter {
	return ctx.rsp
}

func (ctx *Context) SetHandlerList(list []handlerFuncEntry) {
	ctx.hEntryList = list
}

func (ctx *Context) SetParamTable(table map[string]string) {
	ctx.paramTable = table
}

func (ctx *Context) NextHandler() error {
	ctx.hEntryIndex++
	if ctx.hEntryIndex < len(ctx.hEntryList) {
		entry := ctx.hEntryList[ctx.hEntryIndex]

		fmt.Printf("--> [Pipeline-%d/%d] forward to [%s]\n",
			ctx.hEntryIndex+1, len(ctx.hEntryList), entry.funName)

		if err := entry.fun(ctx); err != nil {
			return err
		}
	}

	return nil
}

var _ IRequest = &Context{}

func (ctx *Context) URI() string {
	return ctx.req.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.req.Method
}

func (ctx *Context) Host() string {
	host := ctx.req.Host
	if host == "" {
		// GET /pub/WWW/TheProject.html HTTP/1.1
		// Host: www.example.org:8080
		// ctx.req.URL.Host = ""
		host = ctx.req.URL.Host
	}
	if host == "" {
		// Host Header will be removed by net/http
		host = ctx.req.Header.Get("Host")
	}

	return host
}

func (ctx *Context) ClientIp() string {
	addr := ctx.req.Header.Get("X-Real-Ip")
	if addr == "" {
		addr = ctx.req.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = ctx.req.RemoteAddr
	}

	return addr
}

func (ctx *Context) Headers() map[string][]string {
	return ctx.req.Header
}

func (ctx *Context) Header(key string) (string, bool) {
	valList := ctx.req.Header.Values(key)
	if len(valList) == 0 {
		return "", false
	}

	return valList[0], true
}

func (ctx *Context) Cookies() map[string][]string {
	ctx.cookieParse.Do(func() {
		if ctx.req != nil {
			// every call for ctx.Cookies() will parse cookies and create a new map only once
			cookies := ctx.req.Cookies()
			for _, cookie := range cookies {
				ctx.cookieTable[cookie.Name] = append(ctx.cookieTable[cookie.Name], cookie.Value)
			}
		}
	})

	return ctx.cookieTable
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookie := ctx.Cookies()
	valList, ok := cookie[key]
	if !ok || len(valList) == 0 {
		return "", false
	}

	return valList[0], true
}

func (ctx *Context) QueryAll() map[string][]string {
	ctx.queryAll.Do(func() {
		if ctx.req != nil {
			// every call for ctx.req.URL.Query() will parse URL and create a new map only once
			ctx.queryTable = ctx.req.URL.Query()
		}
	})

	return ctx.queryTable
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.Atoi(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseInt(valList[0], 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseBool(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList[0], true
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList, true
}

func (ctx *Context) Query(key string) interface{} {
	query := ctx.QueryAll()
	valList, ok := query[key]
	if !ok {
		return nil
	}

	return valList
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	val, err := strconv.Atoi(str.(string))
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	val, err := strconv.ParseInt(str.(string), 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	val, err := strconv.ParseFloat(str.(string), 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	val, err := strconv.ParseFloat(str.(string), 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	val, err := strconv.ParseBool(str.(string))
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	str := ctx.Param(key)
	if str == nil {
		return def, false
	}

	return str.(string), true
}

func (ctx *Context) Param(key string) interface{} {
	if len(ctx.paramTable) == 0 {
		return nil
	}

	val, ok := ctx.paramTable[key]
	if !ok {
		return nil
	}

	return val
}

func (ctx *Context) FormAll(key string) map[string][]string {
	if ctx.req != nil {
		ctx.req.ParseForm()     // only parse once
		return ctx.req.PostForm // form in http request body
	}

	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.Atoi(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseInt(valList[0], 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	val, err := strconv.ParseBool(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) FormString(key string, def string) (string, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList[0], true
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok || len(valList) == 0 {
		return def, false
	}

	return valList, true
}

// TODO: add a function to control max memory for form-base data
// TODO: add a api to control where the file be store
const defaultMultipartMemory = 32 << 20 // 32 MB
func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.req.MultipartForm == nil {
		// call ParseMultipartForm() to control the maxMemory
		if err := ctx.req.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}

	file, fileHeader, err := ctx.req.FormFile(key)
	if err != nil {
		return nil, err
	}

	_ = file.Close() // only return FileHeader, but not File
	return fileHeader, err
}

func (ctx *Context) Form(key string) interface{} {
	form := ctx.FormAll(key)
	valList, ok := form[key]
	if !ok {
		return nil
	}

	return valList
}

// BindJSON would parse request body by JSON.
// If obj is nil, BindJSON() not unmarshal, and only store raw JSON into ctx.req.Body
func (ctx *Context) BindJSON(obj interface{}) error {
	if ctx.req == nil {
		return errors.New("Context's request is empty")
	}

	// get all raw data from ctx.req.Body, may IO block
	body, err := ioutil.ReadAll(ctx.req.Body)
	if err != nil {
		// LimitReader wrapped in net/http may return error
		return err
	}

	// store the raw JSON data back to ctx.req.Body
	// after this method function, we can access raw http request JSON body via ctx.Request().Body
	ctx.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if obj == nil {
		// not unmarshal JSON into a struct object
		return nil
	}

	err = json.Unmarshal(body, obj) // unmarshal body into obj interface
	if err != nil {
		return err
	}

	return nil
}

// BindXML would parse request body by XML.
// If obj is nil, BindXML() not unmarshal, and only store raw XML into ctx.req.Body
func (ctx *Context) BindXML(obj interface{}) error {
	if ctx.req == nil {
		return errors.New("Context's request is empty")
	}

	// get all raw data from ctx.req.Body, may IO block
	body, err := ioutil.ReadAll(ctx.req.Body)
	if err != nil {
		// LimitReader wrapped in net/http may return error
		return err
	}

	// store the raw XML data back to ctx.req.Body
	// after this method function, we can access raw http request XML body via ctx.Request().Body
	ctx.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if obj == nil {
		// not unmarshal XML into a struct object
		fmt.Println("not to parse into boj")
		return nil
	}

	err = xml.Unmarshal(body, obj) // unmarshal body into obj interface
	if err != nil {
		return err
	}

	return nil
}

// TODO: now, call RawBody() will perform a deep copy. return same raw data when second call, avoid deep copy
func (ctx *Context) RawBody() ([]byte, error) {
	if ctx.req == nil {
		return nil, errors.New("Context's request is empty")
	}

	body, err := ioutil.ReadAll(ctx.req.Body)
	if err != nil {
		return nil, err
	}
	ctx.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

/* Context impl IResponse */
var _ IResponse = &Context{}

// append Header
func (ctx *Context) AddHeader(key string, val string) IResponse {
	ctx.rsp.Header().Add(key, val)
	return ctx
}

// delete Header
func (ctx *Context) DelHeader(key string) IResponse {
	ctx.rsp.Header().Del(key)
	return ctx
}

// TODO: ctx.SetCookie("key", "val").MaxAge(10).Domain("xxxx") will be better
// TODO: add a simple method may be better, below method has too much param has to provide
func (ctx *Context) SetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.rsp, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})

	return ctx
}

// avoid below bug:
// 1. call SetStatus()
// 2. call Text(), but the Content-Type can not set to ctx.rsp.cw.header
// only call writeStatus() when fill body, or end of all handler function
func writeStatus(ctx *Context, statusCode int) {
	ctx.rsp.WriteHeader(statusCode)
	ctx.wroteStatus = true
}

func (ctx *Context) SetStatus(statusCode int) IResponse {
	ctx.statusCode = statusCode
	return ctx
}

func (ctx *Context) SetStatusOK() IResponse {
	ctx.statusCode = http.StatusOK
	return ctx
}

func (ctx *Context) Redirect(url string) IResponse {
	ctx.SetStatus(http.StatusMovedPermanently)
	ctx.wroteStatus = true
	http.Redirect(ctx.rsp, ctx.req, url, ctx.statusCode)
	return ctx
}

func (ctx *Context) Text(format string, values ...interface{}) IResponse {
	// set Header and write status code
	ctx.AddHeader("Content-Type", "application/text")
	writeStatus(ctx, ctx.statusCode)

	ctx.rsp.Write([]byte(fmt.Sprintf(format, values...))) // TODO: handle error
	return ctx
}

func (ctx *Context) JSON(obj interface{}) IResponse {
	data, err := json.Marshal(obj)
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		writeStatus(ctx, ctx.statusCode)
		return ctx
	}

	ctx.AddHeader("Content-Type", "application/json")
	writeStatus(ctx, ctx.statusCode)

	ctx.rsp.Write(data)
	return ctx
}

func (ctx *Context) XML(obj interface{}) IResponse {
	data, err := xml.Marshal(obj)
	if err != nil {
		ctx.SetStatus(http.StatusInternalServerError)
		writeStatus(ctx, ctx.statusCode)
		return ctx
	}

	ctx.AddHeader("Content-Type", "application/xml")
	writeStatus(ctx, ctx.statusCode)

	ctx.rsp.Write(data)
	return ctx
}

func (ctx *Context) HtmlFiles(obj interface{}, file ...string) IResponse {
	tmpl := template.Must(template.ParseFiles(file...))
	ctx.AddHeader("Content-Type", "application/html")
	if err := tmpl.Execute(ctx.rsp, obj); err != nil {
		ctx.DelHeader("Content-Type")
		return ctx
	}

	return ctx
}

func (ctx *Context) RawData(data []byte) IResponse {
	// write status code
	writeStatus(ctx, ctx.statusCode)

	ctx.rsp.Write(data) // FIXME: handle error return
	return ctx
}
