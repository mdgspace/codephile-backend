package main

import (
	"github.com/astaxie/beego"
	sentryhttp "github.com/getsentry/sentry-go/http"
	_ "github.com/mdg-iitr/Codephile/conf"
	_ "github.com/mdg-iitr/Codephile/routers"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	beego.RunWithMiddleWares("", sentryHandler.Handle)
}
