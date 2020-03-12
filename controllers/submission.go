package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"log"
	"net/http"
	"time"
)

type SubmissionController struct {
	beego.Controller
}

// @Title All submissions
// @Description Get all submissions of a user(logged-in if uid is empty) across various platforms
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Success 200 {object} []types.Submission
// @Failure 400 invalid uid
// @Failure 404 User/Submission not found
// @router /all [get]
// @router /all/:uid [get]
func (s *SubmissionController) GetAllSubmissions() {
	uidString := s.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = s.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid UID")
		s.ServeJSON()
		return
	}
	subs, err := models.GetAllSubmissions(uid)
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Data["json"] = NotFoundError("User/Submission not found")
		s.ServeJSON()
		return
	} else {
		s.Data["json"] = subs
	}
	_ = s.Ctx.Output.JSON(s.Data["json"], false, false)
}

// @Title Get Submissions
// @Description Get paginated submissions(100 per page) of user(logged-in if uid is empty) across various platforms
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Param	before		query 	string	true  "Time before which submissions to be returned, uses current time if empty or not present"
// @Success 200 {object} []types.Submission
// @Failure 400 invalid uid
// @Failure 404 User/Submission not found
// @router / [get]
// @router /:uid [get]
func (s *SubmissionController) PaginatedSubmissions() {
	uidString := s.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = s.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid UID")
		s.ServeJSON()
		return
	}
	before, err := s.GetInt64("before", time.Now().UTC().Unix())
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid query param value")
		s.ServeJSON()
		return
	}
	if before == 0 {
		before = time.Now().UTC().Unix()
	}
	feed, err := models.GetSubmissions(uid, time.Unix(before, 0))
	if err == mgo.ErrNotFound {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Data["json"] = NotFoundError("User not found")
		s.ServeJSON()
		return
	} else if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		s.Data["json"] = InternalServerError("Internal server error")
		s.ServeJSON()
		return
	} else {
		s.Data["json"] = feed
		s.ServeJSON()
	}
}

// @Title Post
// @Description Triggers saving of user's submissions across a particular platform into database
// @Security token_auth write:submission
// @Param	site		path 	string	true		"Platform site name"
// @Success 200 submission successfully saved
// @Failure 400 site invalid
// @Failure 404 user/handle found
// @Failure 500 server_error
// @router /:site [post]
func (s *SubmissionController) SaveSubmission() {
	uid := s.Ctx.Input.GetData("uid").(bson.ObjectId)
	site := s.GetString(":site")
	if !IsSiteValid(site) {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid contest site")
		s.ServeJSON()
		return
	}

	err := models.AddSubmissions(uid, site)
	if err == UserNotFoundError || err == HandleNotFoundError {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = NotFoundError("User/Handle not found")
		s.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Data["json"] = InternalServerError("Internal server error")
		s.ServeJSON()
		return
	}

	s.Data["json"] = map[string]string{"status": "submission successfully saved"}
	s.ServeJSON()
}

// @Title Filter
// @Description Filter submissions of user on the basis of status, site and tags
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Param	site		path 	string	true		"Website name"
// @Param	status		query 	string	false		"Submission status"
// @Param	tag 		query	string	false		"Submission tag"
// @Success 200 {object} types.CodechefSubmission
// @Failure 400 user not exist
// @Failure 500 server_error
// @router /:site/filter [get]
// @router /:site/:uid/filter [get]
func (s *SubmissionController) FilterSubmission() {
	uidString := s.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = s.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid UID")
		s.ServeJSON()
		return
	}
	status := s.GetString("status")
	site := s.GetString(":site")
	tag := s.GetString("tag")
	subs, err := models.FilterSubmission(uid, status, tag, site)
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Data["json"] = InternalServerError("Internal server error")
		s.ServeJSON()
		return
	}
	s.Data["json"] = subs
	s.ServeJSON()
}
