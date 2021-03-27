package main

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	sentryhttp "github.com/getsentry/sentry-go/http"
	_ "github.com/mdg-iitr/Codephile/conf"
	_ "github.com/mdg-iitr/Codephile/routers"
)

func main() {
	if os.Getenv("PORT") == "8080" {
		fmt.Println("dev mode")
		beego.BConfig.RunMode = "dev"
	} else {
		fmt.Println("prod mode")
		beego.BConfig.RunMode = "prod"
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	beego.RunWithMiddleWares("", sentryHandler.Handle)
}
