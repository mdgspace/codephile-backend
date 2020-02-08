package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
)

type FeedController struct {
	beego.Controller
}

// @Title ContestsFeed
// @Description Provides Data for contests in the Feed
// @Security token_auth read:feed
// @Success 200 {object} models.types.S
// @Failure 403 Error fetching contests
// @router /contests [get]
func (f *FeedController) ContestsFeed(){
	 contests, err := models.ReturnFeedContests()
	 if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(403)
		f.Data["json"] = map[string]string{"status": err.Error()}
		f.ServeJSON()
	 }
	 f.Data["json"] = contests
	 f.ServeJSON()
}

// @Title FriendFeed
// @Description Provides Data for Friend Activity in the Feed
// @Security token_auth read:feed
// @Success 200 {object} models.types.FeedObject
// @Failure 403 Invalid uid
// @router /friend-activity [get]
func (f *FeedController) FriendsFeed() {
	  uid := f.Ctx.Input.GetData("uid").(bson.ObjectId)
		  feed,err := models.ReturnFeedFriends(uid)
          if err == models.ErrGeneric {
			  //feed is altered (inform front-end)
			  f.Data["json"] = feed
			  f.ServeJSON()
		  } else if err != nil {
			  f.Ctx.ResponseWriter.WriteHeader(403)
			  f.Data["json"] = err.Error()
			  f.ServeJSON()
		  } else {
			  f.Data["json"] = feed
			  f.ServeJSON()
		  }
}

