package controllers

import (
	"log"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
)

type FriendsController struct {
	beego.Controller
}

// @Title FollowUser
// @Description Adds the Following user's uid to the database
// @Security token_auth write:follow
// @Param	uid2		query 	string	true  "uid of user to follow"
// @Success 200  {string} user followed
// @Failure 400 bad uid
// @Failure 500 server_error
// @router /follow [post]
func (f *FriendsController) FollowUser() {
	uid1 := f.Ctx.Input.GetData("uid").(primitive.ObjectID)
	uid2 := f.GetString("uid2")
	if uid2 == "" || !primitive.IsValidObjectID(uid2) {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		f.Data["json"] = errors.BadInputError("Invalid UID")
		f.ServeJSON()
		return
	}
	id2, _ := primitive.ObjectIDFromHex(uid2)
	err := models.FollowUser(uid1, id2)
	if err != nil {
		hub := sentry.GetHubFromContext(f.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	}
	//user2 has been followed
	f.Data["json"] = map[string]string{"status": "User Followed"}
	f.ServeJSON()
}

// @Title Un-follow User
// @Description Un-follows the user with the given uid
// @Security token_auth write:follow
// @Param	uid2		query 	string	true  "uid of user to un-follow"
// @Success 200  {string} user un-followed
// @Failure 400 bad uid
// @Failure 500 server_error
// @router /unfollow [post]
func (f *FriendsController) UnFollowUser() {
	userUID := f.Ctx.Input.GetData("uid").(primitive.ObjectID)
	uid2 := f.GetString("uid2")
	if uid2 == "" || !primitive.IsValidObjectID(uid2) {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		f.Data["json"] = errors.BadInputError("Invalid UID")
		f.ServeJSON()
		return
	}
	id2, _ := primitive.ObjectIDFromHex(uid2)
	err := models.UnFollowUser(userUID, id2)
	if err == mongo.ErrNoDocuments {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		f.Data["json"] = errors.NotFoundError("user not found")
		f.ServeJSON()
		return
	} else if err != nil {
		hub := sentry.GetHubFromContext(f.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	}
	f.Data["json"] = map[string]string{"status": "User Un-Followed"}
	f.ServeJSON()
}

// @Title CompareUser
// @Description Compares the data of two users
// @Security token_auth read:follow
// @Param	uid2		query 	string	true  "uid of following"
// @Success 200 {object} types.AllWorldRanks
// @Failure 400 bad uid
// @Failure 500 server_error
// @router /compare [get]
func (f *FriendsController) CompareUser() {
	uid1 := f.Ctx.Input.GetData("uid").(primitive.ObjectID)
	uid2 := f.GetString("uid2")
	if uid2 == "" || !primitive.IsValidObjectID(uid2) {
		f.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		f.Data["json"] = errors.BadInputError("Invalid UID")
		f.ServeJSON()
		return
	}
	id2, _ := primitive.ObjectIDFromHex(uid2)
	worldRanks, err := models.CompareUser(uid1, id2)
	if err != nil {
		hub := sentry.GetHubFromContext(f.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	}
	f.Data["json"] = worldRanks
	f.ServeJSON()
}

// @Title GetFollowing
// @Description Fetches the users the user is following
// @Security token_auth read:follow
// @Success 200 {object} []types.FollowingUser
// @Failure 500 server_error
// @router /following [get]
func (f *FriendsController) GetFollowing() {
	uid := f.Ctx.Input.GetData("uid").(primitive.ObjectID)
	following, err := models.GetFollowingUsers(uid)
	if err != nil {
		hub := sentry.GetHubFromContext(f.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		f.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		f.Data["json"] = errors.InternalServerError("Internal server error")
		f.ServeJSON()
		return
	}
	f.Data["json"] = following
	f.ServeJSON()
}
