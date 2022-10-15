/* put this demo in a app command start */

package appDemo

import (
	"MarineSnow/framework/gin"
	"MarineSnow/framework/log"
	"fmt"
	"net/http"
)

// curl -i http://127.0.0.1:80/log/Print
func logPrintf(ctx *gin.Context) {
	log.Print("log.Print(): ", log.GetLevel().String())
	log.Printf("log.Printf(): with level[%s]\n", log.GetLevel().String())
	log.Println("log.Println(): with level[", log.GetLevel().String(), "]")
	ctx.ISetStatusOK()
}

// curl -i http://127.0.0.1:80/log/Log
func logLog(ctx *gin.Context) {
	log.Log(log.ErrorLevel, "log.Log(): with level[", string(log.ParseLevel(log.ErrorLevel)), "]")
	log.Logf(log.ErrorLevel, "log.Logf(): with level[%s]", string(log.ParseLevel(log.ErrorLevel)))
	log.Logln(log.ErrorLevel, "log.Logln(): with level[", string(log.ParseLevel(log.ErrorLevel)), "]")
	ctx.ISetStatusOK()
}

// curl -i http://127.0.0.1:80/log/Level
func logLevel(ctx *gin.Context) {
	log.SetLevel(log.InfoLevel)
	log.Debugf("in debug log, [%s]", log.GetLevel().String())
	log.Infof("in info log, [%s]", log.GetLevel().String())
	ctx.ISetStatusOK()
}

// curl -i http://127.0.0.1:80/log/WithFields
func logWithFields(ctx *gin.Context) {
	fl1 := log.WithFields(log.Fields{
		"demo-key-1": "demo-string-1",
	})

	fl1.Debugf("fieldLogger 1 message")
	fl2 := fl1.WithFields(log.Fields{
		"demo-key-2": "demo-string-2",
	})

	fl2.Infof("fieldLogger 2 message")
	ctx.ISetStatusOK()
}

func StartAppDemo(ipPort string) {
	fmt.Printf("welcome to App-Demo, start now. Listen on [%s]\n", ipPort)
	core := gin.New()
	server := &http.Server{
		Addr:    ipPort,
		Handler: core,
	}

	/* register HTTP handler and route */
	// request demo
	core.GET("/log/Print", logPrintf)
	core.GET("/log/Log", logLog)
	core.GET("/log/Level", logLevel)
	core.GET("/log/WithFields", logWithFields)

	_ = server.ListenAndServe()

	fmt.Println("bye~, exit App-Demo now.")
}
