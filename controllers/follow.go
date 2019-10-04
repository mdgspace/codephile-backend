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

// @Title CompareUser
// @Description Compares the data of two users
// @Param	uid1		query 	string	true  "uid of follower"
// @Param	uid2		query 	string	true  "uid of following"
// @Success 200 {object} //Enter object type
// @Failure 403 Invalid uid
// @router /compare  [get]
func (f *FollowController) CompareUser(){
	uid1 := f.GetString("uid1")
	uid2 := f.GetString("uid2")
	worldRanks ,err := models.CompareUser(uid1, uid2)      //change assignments
    if err == nil{
		//data has been fetched
		f.Data["json"] = worldRanks
	} else {
		//error
		f.Data["json"] = map[string]string{"status": "error"}
	}
	f.ServeJSON()
}