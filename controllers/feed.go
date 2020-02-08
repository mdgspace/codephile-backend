package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"log"
	"net/http"
)

type FeedController struct {
	beego.Controller
}

// @Title ContestsFeed
// @Description Provides Data for contests in the Feed
// @Security token_auth read:feed
// @Success 200 {object} types.S
// @Failure 500 server_error
// @router /contests [get]
func (f *FeedController) ContestsFeed() {
	contests, err := models.ReturnFeedContests()
	if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	}
	f.Data["json"] = contests
	f.ServeJSON()
}

// @Title FriendFeed
// @Description Provides Data for Friend Activity in the Feed
// @Security token_auth read:feed
// @Success 200 {object} types.FeedObject
// @Failure 500 server_error
// @router /friend-activity [get]
func (f *FeedController) FriendsFeed() {
	uid := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	feed, err := models.ReturnFeedFriends(uid)
	if err == models.ErrGeneric {
		//feed is altered (inform front-end)
		f.Data["json"] = feed
		f.ServeJSON()
	} else if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	} else {
		f.Data["json"] = feed
		f.ServeJSON()
	}
}
