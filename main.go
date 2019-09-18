package main

import (
	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/routers"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	beego.Run()
}
