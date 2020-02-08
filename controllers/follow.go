package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
)

type FollowController struct {
	beego.Controller
}

// @Title FollowUser
// @Description Adds the Following user's uid to the database
// @Security token_auth write:follow
// @Param	uid2		query 	string	true  "uid of following"
// @Success 200  {string} user followed
// @Failure 403 Invalid uid
// @router /following [post]
func (f *FollowController) FollowUser(){
	uid1 := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	uid2 := f.GetString("uid2")
	err := models.FollowUser(uid1, uid2)
    if err == nil{
		//user2 has been followed
		f.Data["json"] = map[string]string{"status":"User Followed"}
	} else {
		f.Ctx.ResponseWriter.WriteHeader(403)
		f.Data["json"] = map[string]string{"status": err.Error()}
	}
	f.ServeJSON()
}

// @Title CompareUser
// @Description Compares the data of two users
// @Security token_auth read:follow
// @Param	uid2		query 	string	true  "uid of following"
// @Success 200 {object} models.types.AllWorldRanks
// @Failure 403 Invalid uid
// @router /compare [get]
func (f *FollowController) CompareUser(){
	uid1 := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	uid2 := f.GetString("uid2")
	worldRanks ,err := models.CompareUser(uid1, uid2)      //change assignments
    if err == nil{
		//data has been fetched
		f.Data["json"] = worldRanks
	} else {
		//error
		f.Ctx.ResponseWriter.WriteHeader(403)
		f.Data["json"] = map[string]string{"status": err.Error()}
	}
	f.ServeJSON()
}


// @Title GetFollowing
// @Description Fetches the users the user is following
// @Security token_auth read:follow
// @Success 200 {object} []models.types.Following
// @Failure 403 Invalid uid
// @router /following [get]
func (f *FollowController) GetFollowing(){
	uid := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	following , err := models.GetFollowingUsers(uid)
    if err == nil{
		//data has been fetched
		f.Data["json"] = following
	} else {
		//error
		f.Ctx.ResponseWriter.WriteHeader(403)
		f.Data["json"] = map[string]string{"status": err.Error()}
	}
	f.ServeJSON()
}