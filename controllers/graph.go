package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
)

type GraphController struct {
	beego.Controller
}

// @Title Activity Graph
// @Description Gives the activity graph for a user with given uid, (Logged-in user if uid is empty)
// @Security token_auth read:user
// @Param	uid		path 	string	false		"uid of user"
// @Success 200 {object} types.ActivityGraph
// @Failure 401 : Unauthorized
// @Failure 400 :uid is invalid
// @Failure 404 user not found
// @Failure 500 server_error
// @router /activity [get]
// @router /activity/:uid [get]
func (g *GraphController) GetActivityGraph() {
	uidString := g.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = g.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		g.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		g.Data["json"] = BadInputError("Invalid UID")
		g.ServeJSON()
		return
	}
	graphData, err := models.GetActivityGraph(uid)
	if err != nil {
		hub := sentry.GetHubFromContext(g.Ctx.Request.Context())
		hub.CaptureException(err)
		g.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		g.Data["json"] = InternalServerError("Server error.. Please report to admin")
		g.ServeJSON()
		return
	}
	g.Data["json"] = graphData
	g.ServeJSON()
}

// @Title Submissions Status
// @Description Gives the count of the various different submissions of the user with a uid (Logged-in user if uid is empty)
// @Security token_auth read:user
// @Param	uid		path 	string	false		"uid of user"
// @Success 200 {object} types.StatusCounts
// @Failure 401 : Unauthorized
// @Failure 400 :uid is invalid
// @Failure 404 user not found
// @Failure 500 server_error
// @router /status [get]
// @router /status/:uid [get]
func (g *GraphController) GetStatusCounts() {
	uidString := g.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = g.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		g.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		g.Data["json"] = BadInputError("Invalid UID")
		g.ServeJSON()
		return
	}
	status, err := models.GetStatusCounts(uid)
	if err != nil {
		hub := sentry.GetHubFromContext(g.Ctx.Request.Context())
		hub.CaptureException(err)
		g.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		g.Data["json"] = InternalServerError("Server error.. Please report to admin")
		g.ServeJSON()
		return
	}
	g.Data["json"] = status
	g.ServeJSON()
}
