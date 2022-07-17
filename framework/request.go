package framework

import "mime/multipart"

type IRequest interface {
	/* basic information */
	URI() string
	Method() string
	Host() string
	ClientIp() string

	/* HTTP Header */
	Headers() map[string][]string
	Header(key string) (string, bool)

	// cookie
	Cookies() map[string][]string
	Cookie(key string) (string, bool)

	/* HTTP Body */
	// Example URL: http://localhost/demo?id=123&name=name-string
	// id and name is key; 123 and name-string is return value
	// default must be set
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// parameters in URL
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	// form-base data
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	// format request body
	BindJSON(obj interface{}) error // JSON body
	BindXML(obj interface{}) error  // XML body
	RawBody() ([]byte, error)       // other format
}
