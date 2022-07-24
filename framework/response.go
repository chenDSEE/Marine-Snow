package framework

type IResponse interface {
	/* Header */
	AddHeader(key string, val string) IResponse
	DelHeader(key string) IResponse
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	/* status line */
	SetStatus(code int) IResponse
	SetStatusOK() IResponse
	Redirect(url string) IResponse

	/* body */
	Text(format string, values ...interface{}) IResponse
	JSON(obj interface{}) IResponse
	XML(obj interface{}) IResponse
	HtmlFiles(obj interface{}, file ...string) IResponse
	//Jsonp(obj interface{}) IResponse
	RawData(data []byte) IResponse // not to set Content-Type
}
