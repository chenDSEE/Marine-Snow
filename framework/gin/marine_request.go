package gin

import (
	"strconv"
)

// FIXME: IRequest to make gin.Context own same ability with MarineSnow.Context.
// But IRequest interface is not a good idea, fix this latter.
// TODO: why still define a interface for gin/Context ? I think it is no need
type IRequest interface {
	/* basic information */
	URI() string
	Method() string
	Host() string
	ClientIP() string // replaced by gin.Context.ClientIP()

	/* HTTP Header */
	GetHeaders() map[string][]string
	//TryGetHeader(key string) (string, bool) // replaced by gin.Context.GetHeader()

	// cookie
	Cookies() map[string][]string
	//Cookie(key string) (string, bool) // replace by gin.Context.Cookie()

	/* HTTP Body */
	// Example URL: http://localhost/demo?id=123&name=name-string
	// id and name is key; 123 and name-string is return value
	// default must be set
	// query in URL
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)
	DefaultQueryRaw(key string) interface{}

	// parameters in URL
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	DefaultParamRaw(key string) interface{}

	// form-base data
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	//DefaultFormFile(key string) (*multipart.FileHeader, error) // impl by gin.Context.FormFile()
	DefaultFormRaw(key string) interface{}

	// format request body
	// replaced by gin.Context.BindJSON(), gin.Context.BindXML()
	//BindJSON(obj interface{}) error // JSON body
	//BindXML(obj interface{}) error  // XML body
	//RawBody() ([]byte, error)       // other format
}

var _ IRequest = &Context{}

/* basic information */
func (ctx *Context) URI() string { // TODO: same with gin.Context.FullPath() ?
	return ctx.Request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) Host() string {
	host := ctx.Request.Host
	if host == "" {
		// GET /pub/WWW/TheProject.html HTTP/1.1
		// Host: www.example.org:8080
		// ctx.req.URL.Host = ""
		host = ctx.Request.URL.Host
	}
	if host == "" {
		// Host Header will be removed by net/http
		host = ctx.Request.Header.Get("Host")
	}

	return host
}

/* HTTP Header */
func (ctx *Context) GetHeaders() map[string][]string {
	return ctx.Request.Header
}

// cookie
func (ctx *Context) initCookiesCache() {
	if ctx.cookieCache == nil {
		if ctx.Request != nil {
			// every call for ctx.Cookies() will parse cookies and create a new map only once
			cookies := ctx.Request.Cookies()
			ctx.cookieCache = make(map[string][]string, len(cookies))
			for _, cookie := range cookies {
				ctx.cookieCache[cookie.Name] = append(ctx.cookieCache[cookie.Name], cookie.Value)
			}
		} else {
			ctx.cookieCache = map[string][]string{}
		}
	}
}

func (ctx *Context) Cookies() map[string][]string {
	ctx.initCookiesCache()
	return ctx.cookieCache
}

// query in URL
func (ctx *Context) DefaultQueryInt(key string, def int) (int, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.Atoi(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseInt(valList[0], 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseBool(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultQueryString(key string, def string) (string, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	return valList[0], true
}

func (ctx *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return def, false
	}

	return valList, true
}

func (ctx *Context) DefaultQueryRaw(key string) interface{} {
	valList, ok := ctx.GetQueryArray(key)
	if !ok {
		return nil
	}

	return valList
}

// parameters in URL
func (ctx *Context) DefaultParamInt(key string, def int) (int, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseBool(str)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultParamString(key string, def string) (string, bool) {
	str, ok := ctx.params.Get(key)
	if !ok {
		return def, false
	}

	return str, true
}

func (ctx *Context) DefaultParamRaw(key string) interface{} {
	str, ok := ctx.params.Get(key)
	if !ok {
		return nil
	}

	return str
}

// form-base data
func (ctx *Context) DefaultFormInt(key string, def int) (int, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.Atoi(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseInt(valList[0], 10, 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 64)
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseFloat(valList[0], 32)
	if err != nil {
		return def, false
	}

	return float32(val), true
}

func (ctx *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	val, err := strconv.ParseBool(valList[0])
	if err != nil {
		return def, false
	}

	return val, true
}

func (ctx *Context) DefaultFormString(key string, def string) (string, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	return valList[0], true
}

func (ctx *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return def, false
	}

	return valList, true
}

func (ctx *Context) DefaultFormRaw(key string) interface{} {
	valList, ok := ctx.GetPostFormArray(key)
	if !ok {
		return nil
	}

	return valList
}
