package main

import (
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/routers"
	"github.com/mdg-iitr/Codephile/services/scheduler"
)

func main() {
	
	//goroutine for the scheduler
	go scheduler.StartScheduling()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/docs"] = "swagger"
	}
	beego.Run()
}
