package controllers

import "github.com/astaxie/beego"

type HomePageController struct {
	beego.Controller
}

func (this *HomePageController) Get() {
	this.TplName = "index.html"
	_ = this.Render()
}
