package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"log"
	"net/http"
	"time"
)

type FeedController struct {
	beego.Controller
}

// @Title ContestsFeed
// @Description Provides Data for contests in the Feed
// @Security token_auth read:feed
// @Success 200 {object} types.Result
// @Failure 500 server_error
// @router /contests [get]
func (f *FeedController) ContestsFeed() {
	contests, err := models.GetContestsFeed()
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

// @Title All Feed
// @Description Gives all submission feed of the user
// @Security token_auth read:feed
// @Success 200 {object} []types.FeedObject
// @Failure 500 server_error
// @router /friend-activity/all [get]
func (f *FeedController) AllFeed() {
	uid := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	feed, err := models.GetAllFeed(uid)
	if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	} else {
		f.Data["json"] = feed
		_ = f.Ctx.Output.JSON(feed, false, false)
	}
}

// @Title FriendFeed
// @Description Gives submission feed in paginated manner giving 100 submissions at a time
// @Security token_auth read:feed
// @Param	before		query 	string	true  "Time before which feed to be returned, uses current time if empty or not present"
// @Success 200 {object} []types.FeedObject
// @Failure 400 invalid before value
// @Failure 500 server_error
// @router /friend-activity [get]
func (f *FeedController) PaginatedFeed() {
	uid := f.Ctx.Input.GetData("uid").(bson.ObjectId)
	before, err := f.GetInt64("before", time.Now().UTC().Unix())
	if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		f.Data["json"] = errors.BadInputError("Invalid query param value")
		f.ServeJSON()
		return
	}
	if before == 0 {
		before = time.Now().UTC().Unix()
	}
	feed, err := models.GetFeed(uid, time.Unix(before, 0))
	if err != nil {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	} else {
		f.Data["json"] = feed
		_ = f.Ctx.Output.JSON(feed, false, false)
	}
}
