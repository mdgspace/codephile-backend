package main

import (
	"log"
	"os"

	"github.com/astaxie/beego"
	sentryhttp "github.com/getsentry/sentry-go/http"
	_ "github.com/mdg-iitr/Codephile/conf"
	_ "github.com/mdg-iitr/Codephile/routers"
)

func main() {
	env := os.Getenv("PORT")
	if env == "8080" {
		log.Println("dev mode")
		beego.BConfig.RunMode = "dev"
	} else {
		log.Println("prod mode")
		beego.BConfig.RunMode = "prod"
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	beego.RunWithMiddleWares("", sentryHandler.Handle)
}
