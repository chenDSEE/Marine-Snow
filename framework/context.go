package framework

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	rsp http.ResponseWriter
	req *http.Request

	ctx context.Context

	hEntryIndex int
	hEntryList  []handlerFuncEntry

	// query data in URL
	queryAll   sync.Once
	queryTable map[string][]string

	// param data in URL
	paramTable map[string]string
}

func NewContext(rsp http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		rsp:         rsp,
		req:         req,
		ctx:         req.Context(),
		hEntryIndex: -1,
		hEntryList:  nil,
		queryTable:  map[string][]string{},
		paramTable:  map[string]string{},
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

func (ctx *Context) QueryAll() map[string][]string {
	ctx.queryAll.Do(func() {
		if ctx.req != nil {
			// every call for ctx.req.URL.Query() will parse URL and create a new map
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
