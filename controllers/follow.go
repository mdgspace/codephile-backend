package controllers

import (
	"github.com/astaxie/beego"
	// "github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/Follow"
)

type FollowController struct {
	beego.Controller
}

// @Title FollowUser
// @Description Adds the "Following" uid to the database
// @Param	uid1		query 	string	true  "uid of follower"
// @Param	uid2		query 	string	true  "uid of following"
// @Success 200 {string} user followed
// @Failure 403 Invalid uid
// @router /following  [post]
func (f *FollowController) FollowUser(){
	uid1 := f.GetString("uid1")
	uid2 := f.GetString("uid2")
	err := models.FollowUser(uid1, uid2)
    if err == nil{
		//user2 has been followed
		f.Data["json"] = map[string]string{"status":"User Followed"}
	} else {
		f.Data["json"] = map[string]string{"status": "error"}
	}
	f.ServeJSON()
}