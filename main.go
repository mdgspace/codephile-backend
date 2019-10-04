package main

import (
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/routers"
	"github.com/mdg-iitr/Codephile/services/Scheduler"
)

func main() {
	scheduler.StartScheduling()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	beego.Run()
}
