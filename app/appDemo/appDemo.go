/* put this demo in a app command start */

package appDemo

import (
	"MarineSnow/framework/config"
	"MarineSnow/framework/config/env"
	"MarineSnow/framework/gin"
	"MarineSnow/provider/demo"
	"fmt"
	"net/http"
)

// curl -i http://127.0.0.1:80/configuration
func dumpAll(ctx *gin.Context) {
	fmt.Printf("===== %s =====\n", env.GetBasicInfo())
	env.DumpAll()
	ctx.ISetStatusOK()
}

// curl -i http://127.0.0.1:80/configuration/pair?key=not-exist
// curl -i http://127.0.0.1:80/configuration/pair?kkk=not-exist
// curl -i http://127.0.0.1:80/configuration/pair?key=ENV_INFO
// curl -i http://127.0.0.1:80/configuration/pair?key=ENV_VERSION
// curl -i http://127.0.0.1:80/configuration/pair?key=HOSTNAME
// ENV_MODE="demoAppTest" ./demoApp app start; curl -i http://127.0.0.1:80/configuration/pair?key=ENV_MODE
func pairQuery(ctx *gin.Context) {
	key, ok := ctx.DefaultQueryString("key", "")
	if !ok || env.IsExist(key) != true {
		ctx.ISetStatus(http.StatusNotFound)
		return
	}

	ctx.IText("%s=%s\n", key, env.Get(key))
	ctx.ISetStatusOK()
}

type appCfg struct {
	Network networkCfg `yaml:"network"`
	Log     logCfg     `yaml:"log"`
	Demo    demoCfg    `yaml:"demo"`
}

type networkCfg struct {
	Host string `yaml:"host"`
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Info string `yaml:"info"`
}

type logCfg struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type demoCfg struct {
	Password  string   `yaml:"password"`
	SliceDemo []string `yaml:"sliceDemo"`
}

func StartAppDemo(ipPort string) {
	fmt.Printf("welcome to App-Demo, start now. Listen on [%s]\n", ipPort)
	core := gin.New()
	server := &http.Server{
		Addr:    ipPort,
		Handler: core,
	}

	cfgDecoder := config.NewDecoder("./config/develop/app.yaml", "yaml")
	cfg := appCfg{}
	err := cfgDecoder.LoadConfig(&cfg)
	if err != nil {
		fmt.Printf("Error to load config:[%s]\n", err.Error())
	}
	fmt.Printf("config:[%+v]\n", cfg)

	// register CounterServiceProvider into Container
	_ = core.Register(&demo.CounterServiceProvider{})

	/* register HTTP handler and route */
	// request demo
	core.GET("/configuration", dumpAll)

	configGroup := core.Group("/configuration")
	configGroup.GET("/pair", pairQuery)

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit App-Demo now.")
}
