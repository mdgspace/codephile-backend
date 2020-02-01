package controllers

import (
	// "encoding/json"
	"github.com/astaxie/beego"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/errors"
	"log"
	"net/http"

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
// @Success 200 {object} models.types.S
// @Failure 500 error
// @router / [get]
func (u *ContestController) GetContests() {
	contests, err := models.ReturnContests()
	if err != nil {
		//handle error
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		u.Data["json"] = errors.InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Data["json"] = contests
	u.ServeJSON()
}

// @Title Get Particular Site's Contests	
// @Description Returns the contests of a specific website
// @Security token_auth read:contests
// @Param	site		path 	string	true		"site name"
// @Success 200 {object} models.types.S
// @Failure 400 incorrect site
// @Failure 500 server_error
// @router /:site [get]
func (u *ContestController) GetSpecificContests() {
	site := u.GetString(":site")
	if !IsSiteValid(site) {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = errors.BadInputError("Invalid contest site")
		u.ServeJSON()
		return
	}
	contests, err := models.ReturnSpecificContests(site)
	if err != nil {
		//handle error
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = errors.InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Data["json"] = contests
	u.ServeJSON()
}
