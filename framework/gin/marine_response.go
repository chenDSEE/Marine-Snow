package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	/* Header */
	IAddHeader(key string, val string) IResponse
	IDelHeader(key string) IResponse
	ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	/* status line */
	ISetStatus(code int) IResponse
	ISetStatusOK() IResponse
	IRedirect(url string) IResponse

	/* body */
	IText(format string, values ...interface{}) IResponse
	IJSON(obj interface{}) IResponse
	IXML(obj interface{}) IResponse
	IHtmlFiles(obj interface{}, file ...string) IResponse
	IJsonp(callbackKey string, obj interface{}) IResponse
	IRawData(data []byte) IResponse // not to set Content-Type
}

// gin.Context.Writer own the method of http.ResponseWriter
// gin.Context.responseWriter = http.ResponseWriter (see in gin.responseWriter.reset())
var _ IResponse = &Context{}

/* Header */
func (ctx *Context) IAddHeader(key string, val string) IResponse {
	ctx.Writer.Header().Add(key, val)
	return ctx
}

func (ctx *Context) IDelHeader(key string) IResponse {
	ctx.Writer.Header().Del(key)
	return ctx
}

// TODO: ctx.SetCookie("key", "val").MaxAge(10).Domain("xxxx") will be better
// TODO: add a simple method may be better, below method has too much param has to provide
func (ctx *Context) ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: ctx.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})

	return ctx
}

/* status line */
// avoid below bug:
// 1. call SetStatus()
// 2. call Text(), but the Content-Type can not set to ctx.rsp.cw.header
// only call writeStatus() when fill body, or end of all handler function
// FIXME: case above
//func writeStatus(ctx *Context, statusCode int) {
//	ctx.rsp.WriteHeader(statusCode)
//	ctx.wroteStatus = true
//}
func (ctx *Context) ISetStatus(code int) IResponse {
	ctx.Writer.WriteHeader(code)
	return ctx
}

func (ctx *Context) ISetStatusOK() IResponse {
	ctx.Writer.WriteHeader(http.StatusOK)
	return ctx
}

func (ctx *Context) IRedirect(url string) IResponse {
	http.Redirect(ctx.Writer, ctx.Request, url, http.StatusMovedPermanently)
	return ctx
}

/* body */
func (ctx *Context) IText(format string, values ...interface{}) IResponse {
	// set Header and write status code
	ctx.IAddHeader("Content-Type", "application/text")
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...))) // TODO: handle error
	return ctx
}

func (ctx *Context) IJSON(obj interface{}) IResponse {
	data, err := json.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}

	ctx.IAddHeader("Content-Type", "application/json")
	ctx.Writer.Write(data)
	return ctx
}

func (ctx *Context) IXML(obj interface{}) IResponse {
	data, err := xml.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}

	ctx.IAddHeader("Content-Type", "application/xml")
	ctx.Writer.Write(data)
	return ctx
}

func (ctx *Context) IHtmlFiles(obj interface{}, file ...string) IResponse {
	tmpl := template.Must(template.ParseFiles(file...))
	ctx.IAddHeader("Content-Type", "application/html")

	if err := tmpl.Execute(ctx.Writer, obj); err != nil {
		ctx.IDelHeader("Content-Type")
		return ctx
	}

	return ctx
}

// TODO: need more investigation
// JSONP request like: https://127.0.0.1:80/jsonp?callbackKey=callbackFunction
func (ctx *Context) IJsonp(callbackKey string, obj interface{}) IResponse {
	// get callback name from URL query
	callbackName, _ := ctx.DefaultQueryString(callbackKey, "callback")
	ctx.IAddHeader("Content-Type", "application/javascript")

	// avoid XSS for web
	callback := template.JSEscapeString(callbackName)

	// format: callback_name(JSON-data)
	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		return ctx
	}

	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		return ctx
	}

	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.Writer.Write(ret)
	if err != nil {
		return ctx
	}

	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx
}

func (ctx *Context) IRawData(data []byte) IResponse {
	ctx.Writer.Write(data) // FIXME: handle error return
	return ctx
}
