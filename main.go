package main

import (
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/routers"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	beego.RunWithMiddleWares("", sentryHandler.Handle)
}
