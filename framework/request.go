package framework

type IRequest interface {
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
	Query(key string) (string, bool)

	// parameters in URL
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) (string, bool)
}
