package controllers

import (
	// "encoding/json"
	"github.com/astaxie/beego"
	//"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
)

//Controller to display contests
type ContestController struct {
	beego.Controller
}

// @Title GetContests
// @Description displays all contests
// @Security token_auth read:contests
// @Success 200 {object} models.S
// @Failure 403 error
// @router / [get]
func (u *ContestController) GetContests() {
	contests := models.ReturnContests()
	u.Data["json"] = contests
	u.ServeJSON()
}

// @Title Get Particular Site's Contests	
// @Description Returns the contests of a specific website
// @Security token_auth read:contests
// @Param	site		path 	string	true		"site name"
// @Success 200 {object} models.S
// @Failure 403 incorrect site or unknown error
// @router /:site [get]
func (u *ContestController) GetSpecificContests() {
	site := u.GetString(":site")
	contests := models.ReturnSpecificContests(site)
	u.Data["json"] = contests
	u.ServeJSON()
}



